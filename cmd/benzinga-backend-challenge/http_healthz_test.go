package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestHealthzHandler(t *testing.T) {
	r := mux.NewRouter()
	h := &healthzHandler{}
	configureHealthzHTTPRoute(r.NewRoute()).Handler(h)
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	w.Flush()
	if w.Code != http.StatusOK {
		t.Log(w.Body)
		t.Fatalf("unexpected code: got %d, want %d", w.Code, http.StatusOK)
	}
}
