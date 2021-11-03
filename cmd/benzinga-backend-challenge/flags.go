package main

import (
	"flag"
	"os"
	"time"
)

type flags struct {
	http    string
	logFile string
	size    int
	delay   time.Duration
}

func getFlags() *flags {
	flg := initFlagsFromEnv()
	fs := flag.NewFlagSet("", flag.ExitOnError)
	fs.StringVar(&flg.http, "http", flg.http, "HTTP ")
	fs.StringVar(&flg.logFile, "log-file", flg.logFile, "log file")
	fs.IntVar(&flg.size, "size", flg.size, "Size")
	fs.DurationVar(&flg.delay, "delay", flg.delay, "Delay")
	_ = fs.Parse(os.Args[1:]) // Ignore error, because it exits on error
	return flg
}

func initFlagsFromEnv() *flags {
	return &flags{
		http: ":8080",
	}
}
