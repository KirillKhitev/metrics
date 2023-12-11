package main

import (
	"flag"
)

var flagAddrRun string

func parseFlags() {
	flag.StringVar(&flagAddrRun, "a", "localhost:8080", "address and port to run server")
	flag.Parse()
}
