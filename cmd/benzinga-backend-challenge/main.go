package main

import (
	"context"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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
	log.WithFields(log.Fields{
		"version": version,
		"app":     appName,
	}).Info("running")
	flg := getFlags()
	ctx := context.Background()
	dic := newDIContainer(flg)
	err := runHTTPServer(ctx, dic, flg.http)
	if err != nil {
		return errors.Wrap(err, "HTTP server")
	}
	return nil
}
