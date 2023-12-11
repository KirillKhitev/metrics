package main

import (
	"flag"
)

var flagAddrRun string
var flagPollInterval int
var flagReportInterval int

func parseFlags() {
	flag.StringVar(&flagAddrRun, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&flagPollInterval, "p", 2, "poll metrics interval")
	flag.IntVar(&flagReportInterval, "r", 10, "send metrics report interval")
	flag.Parse()
}
