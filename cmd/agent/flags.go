package main

import (
	"flag"
	"time"
)

var flagAddrRun string
var flagPollInterval time.Duration
var flagReportInterval time.Duration

func parseFlags() {
	flag.StringVar(&flagAddrRun, "a", "localhost:8080", "address and port to run server")
	flag.DurationVar(&flagPollInterval, "p", 2*1000000000, "poll metrics interval")
	flag.DurationVar(&flagReportInterval, "r", 10*1000000000, "send metrics report interval")
	flag.Parse()
}
