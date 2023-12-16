package main

import (
	"flag"
	"os"
)

var flagAddrRun string

func parseFlags() {
	flag.StringVar(&flagAddrRun, "a", "localhost:8080", "address and port to run server")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagAddrRun = envRunAddr
	}
}
