package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type webhookForwarder struct {
	endpoint      string
	batchSize     int
	batchInterval time.Duration
}

func newWebhookForwarderHandler(dic *diContainer) *webhookForwarder {
	return &webhookForwarder{
		endpoint:      dic.flags.postEndpoint,
		batchSize:     dic.flags.batchSize,
		batchInterval: dic.flags.batchInterval,
	}
}

func newWebhookForwarderDIProvider(dic *diContainer) func() *webhookForwarder {
	var w *webhookForwarder
	var mu sync.Mutex
	return func() *webhookForwarder {
		mu.Lock()
		defer mu.Unlock()
		if w == nil {
			w = newWebhookForwarderHandler(dic)
		}
		return w
	}
}

func (w *webhookForwarder) forward(ctx context.Context, msgStream <-chan *logHTTPHandlerRequestBody, errCh chan<- error) {
	eventsPayload := []*logHTTPHandlerRequestBody{}
	deadline := time.After(w.batchInterval)
	for msg := range msgStream {
		eventsPayload = append(eventsPayload, msg)
		// if batch size was set
		if w.batchSize > 0 && len(eventsPayload) >= w.batchSize {
			w.forwardEvents(
				ctx,
				eventsPayload,
				errCh,
			)
			// clear cache
			eventsPayload = nil
			// reset deadline
			deadline = time.After(w.batchInterval)
		} else if w.batchInterval > 0 { // if batchInterval was set, try to check if its reached
			select {
			case <-deadline:
				w.forwardEvents(
					ctx,
					eventsPayload,
					errCh,
				)
				// reset deadline
				deadline = time.After(w.batchInterval)
				// clear cache
				eventsPayload = nil
			default:
				// in case deadline is not reached, continue
				continue
			}
		}
	}
}

func (w *webhookForwarder) forwardEvents(ctx context.Context, eventsPayload []*logHTTPHandlerRequestBody, errCh chan<- error) {
	// set time was probably reached, however no new payload was received from /log
	if len(eventsPayload) == 0 {
		return
	}
	timeStart := time.Now()
	statusCode, err := w.forwardWithRetries(
		ctx,
		eventsPayload,
	)
	if err != nil {
		err = errors.Wrap(err, "forward with retries exhausted")
		errCh <- err
		return
	}
	log.WithFields(
		log.Fields{
			"latency":          time.Since(timeStart),
			"http_status_code": statusCode,
			"batch_size":       len(eventsPayload),
		},
	).Info("request success")
}

func (w *webhookForwarder) forwardWithRetries(ctx context.Context, eventsPayload []*logHTTPHandlerRequestBody) (int, error) {
	// Retrying won't help as body is malformed
	bodyWebhook, err := json.Marshal(eventsPayload)
	if err != nil {
		return 0, errors.Wrap(err, "marshal")
	}
	// Retrying won't help as its an issue with url parse
	req, err := http.NewRequest(
		http.MethodPost,
		w.endpoint,
		bytes.NewBuffer(bodyWebhook),
	)
	if err != nil {
		return 0, errors.Wrap(err, "new HTTP request")
	}
	req.Header.Add("Content-Type", "application/json")
	req = req.WithContext(ctx)
	retries := 0
	var lastErr error
	for {
		// return if retires exceeds 3 times and one original try
		if retries > 3 {
			return 0, lastErr
		}
		// sleep before each retry but not first try
		if retries >= 1 {
			time.Sleep(2 * time.Second)
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			err = errors.Wrap(err, "DO http client request")
			lastErr = err
			retries++
		}
		defer res.Body.Close() //nolint:errcheck
		if res.StatusCode >= 200 && res.StatusCode < 300 {
			return res.StatusCode, nil
		}
	}
}
