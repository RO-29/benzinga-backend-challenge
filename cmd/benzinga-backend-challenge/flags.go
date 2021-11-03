package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"
)

type flags struct {
	http          string
	logFile       string
	batchSize     int
	batchInterval time.Duration
	postEndpoint  string
}

func getFlags() *flags {
	flg := &flags{
		http:          ":8080",
		batchSize:     10,
		batchInterval: 10 * time.Second,
	}
	initFlagsFromEnv(flg)
	fs := flag.NewFlagSet("", flag.ExitOnError)
	fs.StringVar(&flg.http, "http", flg.http, "HTTP ")
	fs.StringVar(&flg.logFile, "log-file", flg.logFile, "log file")
	fs.IntVar(&flg.batchSize, "batch-size", flg.batchSize, "Batch Size")
	fs.DurationVar(&flg.batchInterval, "batch-interval", flg.batchInterval, "Batch Interval")
	fs.StringVar(&flg.postEndpoint, "post-endpoint", flg.logFile, "Post Endpoint")
	_ = fs.Parse(os.Args[1:]) // Ignore error, because it exits on error
	checkFlags(flg)
	return flg
}

func checkFlags(flg *flags) {
	if flg.postEndpoint == "" {
		log.Fatal("mandatory post endpoint not set")
	}
	_, err := url.Parse(flg.postEndpoint)
	if err != nil {
		log.Fatalf("\ninvalid post endpoint: '%s'", flg.postEndpoint)
	}
	if flg.batchInterval == 0 && flg.batchSize == 0 {
		log.Fatal("either of batch size or batch interval must be set")
	}
}

func initFlagsFromEnv(flg *flags) {
	setPostEndpointEnv(flg)
	setBatchSizeEnv(flg)
	setBatchIntervalEnv(flg)
}

const (
	postEndPointEnvVar  = "WEBHOOK_POST_ENDPOINT"
	batchSizeEnvVar     = "WEBHOOK_BATCH_SIZE"
	batchIntervalEnvVar = "WEBHOOK_BATCH_INTERVAL"
)

func setPostEndpointEnv(flg *flags) {
	u, ok := os.LookupEnv(postEndPointEnvVar)
	if !ok {
		log.Printf("no valid post endpoint found in env: %s\n", postEndPointEnvVar)
		return
	}
	_, err := url.Parse(u)
	if err != nil {
		log.Printf("invalid post endpoint: '%s' found in env: %s\n", u, postEndPointEnvVar)
		return
	}
	flg.postEndpoint = u
}

func setBatchSizeEnv(flg *flags) {
	bs, ok := os.LookupEnv(batchSizeEnvVar)
	if !ok {
		log.Printf("no valid batch size found in env: %s\n", batchSizeEnvVar)
		return
	}
	bsI, err := strconv.Atoi(bs)
	if err != nil {
		log.Printf("no valid batch size found in env: %s\n", batchSizeEnvVar)
		return

	}
	flg.batchSize = bsI
}

func setBatchIntervalEnv(flg *flags) {
	bi, ok := os.LookupEnv(batchIntervalEnvVar)
	if !ok {
		log.Printf("no valid batch interval found in env: %s\n", batchIntervalEnvVar)
		return
	}
	biT, err := time.ParseDuration(bi)
	if err != nil {
		log.Printf("no valid batch interval found in env: %s\n", batchIntervalEnvVar)
		return
	}
	flg.batchInterval = biT
}
