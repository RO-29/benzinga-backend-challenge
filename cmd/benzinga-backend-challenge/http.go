package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func runHTTPServer(ctx context.Context, dic *diContainer, addr string) error {
	h, err := dic.httpHandler()
	if err != nil {
		return errors.Wrap(err, "get http handler")
	}
	log.WithFields(log.Fields{
		"addr": addr,
	}).Info("Start HTTP server on")
	srv := &http.Server{
		Addr:    addr,
		Handler: h,
	}
	go runForwarderConsumer(ctx, dic)
	err = srv.ListenAndServe()
	if err != nil {
		return errors.Wrap(err, "listen and serve")
	}

	log.Info("Stopped HTTP server")
	return nil
}

var logBufferMsgs chan *logHTTPHandlerRequestBody

func runForwarderConsumer(ctx context.Context, dic *diContainer) {
	bufferSize := 100
	if dic.flags.batchSize > 0 {
		bufferSize = dic.flags.batchSize * 2
	}
	//bufferSize twice the batch size in case processing slow and we don't want /log to wait to send msg to buffered channel, *100* if batchSize was  not set.
	/*
		to process faster, we an start runForwarderConsumer in waitgroup mode with adding more consumer as we want but throughput must be higher to add more consumers, starting with one consumer for now
	*/
	logBufferMsgs = make(chan *logHTTPHandlerRequestBody, bufferSize)
	fw := dic.webhookForwarder()
	errCh := make(chan error)
	go fw.forward(ctx, logBufferMsgs, errCh)
	// wait for err signal from forwarder, exit in case of error
	err := <-errCh
	if err != nil {
		log.WithField("err: ", err).Fatal("exit")
	}
}

func newHTTPHandler(dic *diContainer) (http.Handler, error) {
	r, err := dic.httpRouter()
	if err != nil {
		return nil, errors.Wrap(err, "router")
	}
	return r, nil
}

func newHTTPHandlerDIProvider(dic *diContainer) func() (http.Handler, error) {
	var v http.Handler
	var mu sync.Mutex
	return func() (_ http.Handler, err error) {
		mu.Lock()
		defer mu.Unlock()
		if v == nil {
			v, err = newHTTPHandler(dic)
		}
		return v, err
	}
}

func newHTTPRouter(dic *diContainer) (*mux.Router, error) {
	r := mux.NewRouter()
	err := registerHTTPHandlers(r, dic.httpHandlers)
	if err != nil {
		return nil, errors.Wrap(err, "register")
	}
	return r, nil
}

func newHTTPRouterDIProvider(dic *diContainer) func() (*mux.Router, error) {
	var r *mux.Router
	var mu sync.Mutex
	return func() (*mux.Router, error) {
		mu.Lock()
		defer mu.Unlock()
		var err error
		if r == nil {
			r, err = newHTTPRouter(dic)
		}
		return r, err
	}
}

func registerHTTPHandlers(r *mux.Router, hs *httpHandlers) error {
	for _, v := range []struct {
		name      string
		configure func(*mux.Route) *mux.Route
		handler   func() (http.Handler, error)
	}{
		{
			name:      "log",
			configure: configureLogHTTPRoute,
			handler:   hs.logHandler,
		},
		{
			name:      "healthz",
			configure: configureHealthzHTTPRoute,
			handler:   hs.healthzHandler,
		},
	} {
		h, err := v.handler()
		if err != nil {
			return errors.Wrap(err, v.name)
		}
		v.configure(r.NewRoute()).Handler(h)
	}
	return nil
}

type httpHandlers struct {
	logHandler     func() (http.Handler, error)
	healthzHandler func() (http.Handler, error)
}

func newHTTPHandlers(dic *diContainer) *httpHandlers {
	return &httpHandlers{
		logHandler:     newLogHandlerDIProvider(dic),
		healthzHandler: newHealthzHandlerDIProvider(),
	}
}

func onHTTPError(_ context.Context, w http.ResponseWriter, _ *http.Request, err error, resp *httpErrorResponse) {
	hd := w.Header()
	hd.Set("Content-Type", "application/json")
	w.WriteHeader(resp.Code)
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	_ = enc.Encode(resp)
	_, _ = w.Write(buf.Bytes())
	log.WithFields(log.Fields{
		"err": err,
	}).Error()
}

type httpErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"-"`
}
