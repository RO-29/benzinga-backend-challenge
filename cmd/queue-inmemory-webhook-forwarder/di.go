package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type diContainer struct {
	flags *flags

	httpHandler      func() (http.Handler, error)
	httpRouter       func() (*mux.Router, error)
	httpHandlers     *httpHandlers
	webhookForwarder func() *webhookForwarder
}

func newDIContainer(flg *flags) *diContainer {
	dic := &diContainer{
		flags: flg,
	}
	dic.httpHandlers = newHTTPHandlers(dic)
	dic.httpRouter = newHTTPRouterDIProvider(dic)
	dic.httpHandler = newHTTPHandlerDIProvider(dic)
	dic.webhookForwarder = newWebhookForwarderDIProvider(dic)
	return dic
}
