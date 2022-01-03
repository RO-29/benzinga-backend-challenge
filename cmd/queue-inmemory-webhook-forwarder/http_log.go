package main

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type logHandler struct{}

func newLogHandler(dic *diContainer) *logHandler {
	return &logHandler{}
}

func newLogHandlerDIProvider(dic *diContainer) func() (http.Handler, error) {
	var l *logHandler
	var mu sync.Mutex
	return func() (http.Handler, error) {
		mu.Lock()
		defer mu.Unlock()
		if l == nil {
			l = newLogHandler(dic)
		}
		return l, nil
	}
}

func configureLogHTTPRoute(r *mux.Route) *mux.Route {
	return r.Methods(http.MethodPost).Path("/log")
}

func (h *logHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.handle(req.Context(), w, req)
}

type logHTTPHandlerRequestBody struct {
	UserID    int64    `json:"user_id"`
	Total     float64  `json:"total"`
	Title     string   `json:"title"`
	Meta      metaInfo `json:"meta"`
	Completed bool     `json:"completed"`
}

type metaInfo struct {
	Logins       []login      `json:"logins"`
	PhoneNumbers phoneNumbers `json:"phone_numbers"`
}

type login struct {
	Time time.Time `json:"time"`
	IP   string    `json:"ip"`
}

type phoneNumbers struct {
	Home   string `json:"home"`
	Mobile string `json:"mobile"`
}

func (h *logHandler) handle(_ context.Context, w http.ResponseWriter, req *http.Request) {
	log.Info("request received for /log")
	body, err := h.decodeRequestBody(req)
	if err != nil {
		onHTTPError(
			req.Context(),
			w,
			req,
			err,
			&httpErrorResponse{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			})
		return
	}
	// Note it could be in blocking state if the buffer is full and forward consumer is processing
	logBufferMsgs <- body

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("accepted"))
}

func (h *logHandler) decodeRequestBody(req *http.Request) (*logHTTPHandlerRequestBody, error) {
	var v *logHTTPHandlerRequestBody
	err := json.NewDecoder(req.Body).Decode(&v)
	if err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	if v.UserID == 0 {
		return nil, errors.New("can not post empty user id")
	}
	return v, nil
}
