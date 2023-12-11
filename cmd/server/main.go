package main

import (
	"github.com/KirillKhitev/metrics/internal/server"
	"log"
	"net/http"
)

func main() {
	parseFlags()

	log.Fatal(http.ListenAndServe(flagAddrRun, server.GetRouter()))
}
