package main

import (
	"github.com/KirillKhitev/metrics/internal/flags"
	"github.com/KirillKhitev/metrics/internal/gzip"
	"github.com/KirillKhitev/metrics/internal/logger"
	"github.com/KirillKhitev/metrics/internal/server"
	"github.com/KirillKhitev/metrics/internal/storage"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	flags.Args.Parse()

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	if err := logger.Initialize("info"); err != nil {
		return err
	}

	var appStorage storage.Repository

	if flags.Args.DBConnectionString != "" {
		appStorage = &storage.DBStorage{}
	} else {
		appStorage = &storage.MemStorage{}
	}

	if err := appStorage.Init(); err != nil {
		return err
	}

	go saveToFile(appStorage)

	logger.Log.Info("Running server", zap.String("address", flags.Args.AddrRun))

	handler := gzip.Middleware(server.GetRouter(appStorage))

	return http.ListenAndServe(flags.Args.AddrRun, logger.RequestLogger(handler))
}

func saveToFile(appStorage storage.Repository) {
	ticker := make(<-chan time.Time)

	if flags.Args.StoreInterval > 0 {
		ticker = time.Tick(time.Second * time.Duration(flags.Args.StoreInterval))
	}

	terminateSignals := make(chan os.Signal, 1)

	signal.Notify(terminateSignals, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-ticker:
			if err := appStorage.TrySaveToFile(); err != nil {
				logger.Log.Error("Error by save metrics to file", zap.Error(err))
			}

		case <-terminateSignals:
			if err := appStorage.TrySaveToFile(); err != nil {
				logger.Log.Error("Error by save metrics to file", zap.Error(err))
			}
			appStorage.Close()
			os.Exit(1)
		}
	}
}
