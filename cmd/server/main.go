package main

import (
	"github.com/KirillKhitev/metrics/internal/logger"
	"github.com/KirillKhitev/metrics/internal/server"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	parseFlags()

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	if err := logger.Initialize("info"); err != nil {
		return err
	}

	logger.Log.Info("Running server", zap.String("address", flagAddrRun))

	return http.ListenAndServe(flagAddrRun, logger.RequestLogger(server.GetRouter()))
}
