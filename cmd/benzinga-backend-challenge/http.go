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

func runForwarderConsumer(ctx context.Context, dic *diContainer) {
	inputBufferMsgs := make(chan *logHTTPHandlerRequestBody, dic.flags.batchSize*2)
	fw := dic.webhookForwarder()
	errCh := make(chan error)
	go fw.forward(ctx, inputBufferMsgs, errCh)
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
