package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestLogHandler(t *testing.T) {
	r := mux.NewRouter()
	h := &logHandler{}
	body := `{
		"user_id": 1,
		"total": 1.65,
		"title": "delectus aut autem",
		"meta": {
			"logins": [{
				"time": "2020-08-08T01:52:50Z",
				"ip": "0.0.0.0"
			}],
			"phone_numbers": {
				"home": "555-1212",
				"mobile": "123-5555"
			}
		},
		"completed": false
	}
	`
	configureLogHTTPRoute(r.NewRoute()).Handler(h)
	req := httptest.NewRequest(http.MethodPost, "/log", strings.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	w.Flush()
	if w.Code != http.StatusCreated {
		t.Log(w.Body)
		t.Fatalf("unexpected code: got %d, want %d", w.Code, http.StatusCreated)
	}
}

func TestLogHandlerSetErr(t *testing.T) {
	r := mux.NewRouter()
	h := &logHandler{}
	body := `{
		"total": 1.65,
		"title": "delectus aut autem",
		"meta": {
			"logins": [{
				"time": "2020-08-08T01:52:50Z",
				"ip": "0.0.0.0"
			}],
			"phone_numbers": {
				"home": "555-1212",
				"mobile": "123-5555"
			}
		},
		"completed": false
	}
	`
	configureLogHTTPRoute(r.NewRoute()).Handler(h)
	req := httptest.NewRequest(http.MethodPost, "/log", strings.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	w.Flush()
	if w.Code != http.StatusBadRequest {
		t.Log(w.Body)
		t.Fatalf("unexpected code: got %d, want %d", w.Code, http.StatusBadRequest)
	}
}
