package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/pkg/errors"
)

func newServerPostWebHookCall(t *testing.T, expectedBody []*logHTTPHandlerRequestBody, expectedStatusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ct := req.Header.Get("Content-Type")
		if ct != "application/json" {
			err := errors.Errorf("unexpected content type: got %q, want %q", ct, "application/json")
			t.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var body []*logHTTPHandlerRequestBody
		err := json.NewDecoder(req.Body).Decode(&body)
		if err != nil {
			// t.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !reflect.DeepEqual(body, expectedBody) {
			http.Error(w, "unexpected body", http.StatusBadRequest)
		}
		w.WriteHeader(expectedStatusCode)
	}))
}

func TestWebHookCall(t *testing.T) {
	for _, tc := range []struct {
		name               string
		data               []*logHTTPHandlerRequestBody
		expectedBody       []*logHTTPHandlerRequestBody
		expectedStatusCode int
		expectedErr        bool
	}{
		{
			name: "OK",
			data: []*logHTTPHandlerRequestBody{
				{
					UserID: 1,
					Total:  10.0,
				},
			},
			expectedBody: []*logHTTPHandlerRequestBody{
				{
					UserID: 1,
					Total:  10.0,
				},
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "badRequest",
			data: []*logHTTPHandlerRequestBody{
				{
					UserID: 1,
					Total:  14.0,
				},
			},
			expectedBody: []*logHTTPHandlerRequestBody{
				{
					UserID: 1,
					Total:  14.0,
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErr:        true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := newServerPostWebHookCall(
				t,
				tc.expectedBody,
				tc.expectedStatusCode,
			)
			defer func() {
				server.Close()
			}()
			fw := &webhookForwarder{
				endpoint:           server.URL,
				retrySleepInterval: 10 * time.Millisecond,
				retryLimit:         3,
			}
			statusCode, err := fw.forwardWithRetries(
				context.Background(),
				tc.data,
			)
			if err != nil && !tc.expectedErr {
				t.Fatal(err)
			}
			if statusCode != tc.expectedStatusCode {
				t.Fatalf("unexpected status code: got:%d want:%d", statusCode, tc.expectedStatusCode)
			}
		})
	}
}
