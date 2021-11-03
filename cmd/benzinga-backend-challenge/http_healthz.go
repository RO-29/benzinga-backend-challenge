package main

import (
	"context"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

type healthzHandler struct { // TODO Have fun here
}

func newHealthzHandler() *healthzHandler {
	return &healthzHandler{}
}

func newHealthzHandlerDIProvider() func() (http.Handler, error) {
	var h *healthzHandler
	var mu sync.Mutex
	return func() (http.Handler, error) {
		mu.Lock()
		defer mu.Unlock()
		if h == nil {
			h = newHealthzHandler()
		}
		return h, nil
	}
}

func configureHealthzHTTPRoute(r *mux.Route) *mux.Route {
	return r.Methods(http.MethodGet).Path("/healthz")
}

func (h *healthzHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.handle(req.Context(), w, req)
}

func (h *healthzHandler) handle(_ context.Context, w http.ResponseWriter, _ *http.Request) {
	// TODO ideas for health check?
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}
