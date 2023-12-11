package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

var flagAddrRun string
var flagPollInterval int
var flagReportInterval int

func parseFlags() {
	flag.StringVar(&flagAddrRun, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&flagPollInterval, "p", 2, "poll metrics interval")
	flag.IntVar(&flagReportInterval, "r", 10, "send metrics report interval")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagAddrRun = envRunAddr
	}

	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		if val, err := strconv.Atoi(envPollInterval); err == nil {
			flagPollInterval = val
		} else {
			log.Println(fmt.Sprintf("wrong value environment POLL_INTERVAL: %s", envPollInterval))
		}
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		if val, err := strconv.Atoi(envReportInterval); err == nil {
			flagReportInterval = val
		} else {
			log.Println(fmt.Sprintf("wrong value environment REPORT_INTERVAL: %s", flagReportInterval))
		}
	}
}
