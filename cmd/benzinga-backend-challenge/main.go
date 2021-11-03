package main

import (
	"context"
	"log"

	"github.com/pkg/errors"
)

var version = "dev"

const (
	appName = "benzinga-backend-challenge"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(errors.Wrap(err, "run"))
	}
}

func run() error {
	log.Printf("running version: %#v, app: %#v", version, appName)
	flg := getFlags()
	ctx := context.Background()
	dic := newDIContainer(flg)
	err := runHTTPServer(ctx, dic, flg.http)
	if err != nil {
		return errors.Wrap(err, "HTTP server")
	}
	return nil
}
