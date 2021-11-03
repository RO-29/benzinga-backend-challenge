package main

import (
	"context"
	"sync"
)

type webhookForwarder struct {
	endpoint string
}

func newWebhookForwarderHandler(dic *diContainer) *webhookForwarder {
	return &webhookForwarder{
		endpoint: dic.flags.postEndpoint,
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

func (w *webhookForwarder) forward(ctx context.Context, msgStream <-chan logHTTPHandlerRequestBody, errCh chan<- error) {
	_, _, _ = ctx, msgStream, errCh
}
