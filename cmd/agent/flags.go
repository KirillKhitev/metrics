package main

import (
	"flag"
	"log"
	"os"
	"strconv"
)

var flags flagsAgent = flagsAgent{}

type flagsAgent struct {
	AddrRun        string
	PollInterval   int
	ReportInterval int
	RateLimit      int
	Key            string
}

func (f *flagsAgent) ParseFlags() {
	flag.StringVar(&f.AddrRun, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&f.PollInterval, "p", 2, "poll metrics interval")
	flag.IntVar(&f.ReportInterval, "r", 10, "send metrics report interval")
	flag.IntVar(&f.RateLimit, "l", 5, "request to server limit")
	flag.StringVar(&f.Key, "k", "", "key for signature request data")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		f.AddrRun = envRunAddr
	}

	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		if val, err := strconv.Atoi(envPollInterval); err == nil {
			f.PollInterval = val
		} else {
			log.Printf("wrong value environment POLL_INTERVAL: %s", envPollInterval)
		}
	}

	if envRateLimit := os.Getenv("RATE_LIMIT"); envRateLimit != "" {
		if val, err := strconv.Atoi(envRateLimit); err == nil {
			f.RateLimit = val
		} else {
			log.Printf("wrong value environment RATE_LIMIT: %s", envRateLimit)
		}
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		if val, err := strconv.Atoi(envReportInterval); err == nil {
			f.ReportInterval = val
		} else {
			log.Printf("wrong value environment REPORT_INTERVAL: %s", envReportInterval)
		}
	}

	if envKey := os.Getenv("KEY"); envKey != "" {
		f.Key = envKey
	}
}
