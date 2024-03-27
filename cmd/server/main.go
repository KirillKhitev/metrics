package main

import (
	"context"
	"github.com/KirillKhitev/metrics/internal/flags"
	"github.com/KirillKhitev/metrics/internal/gzip"
	"github.com/KirillKhitev/metrics/internal/logger"
	"github.com/KirillKhitev/metrics/internal/server"
	"github.com/KirillKhitev/metrics/internal/signature"
	"github.com/KirillKhitev/metrics/internal/storage"
	"go.uber.org/zap"
	"net/http"
	_ "net/http/pprof"
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

	if err := appStorage.Init(context.Background()); err != nil {
		return err
	}

	go intervalSaveToFile(appStorage)
	go startServer(appStorage)
	go startServerPprof()

	return catchTerminateSignal(appStorage)
}

func startServer(appStorage storage.Repository) error {
	logger.Log.Info("Running server", zap.String("address", flags.Args.AddrRun))

	handler := gzip.Middleware(server.GetRouter(appStorage))
	handler = signature.Middleware(handler)

	return http.ListenAndServe(flags.Args.AddrRun, logger.RequestLogger(handler))
}

func intervalSaveToFile(appStorage storage.Repository) {
	ticker := make(<-chan time.Time)

	if flags.Args.StoreInterval > 0 {
		ticker = time.Tick(time.Second * time.Duration(flags.Args.StoreInterval))
	}

	for {
		<-ticker
		if err := appStorage.TrySaveToFile(); err != nil {
			logger.Log.Error("Error by save metrics to file", zap.Error(err))
		}
	}
}

func catchTerminateSignal(appStorage storage.Repository) error {
	terminateSignals := make(chan os.Signal, 1)

	signal.Notify(terminateSignals, syscall.SIGINT, syscall.SIGTERM)

	<-terminateSignals
	if err := appStorage.TrySaveToFile(); err != nil {
		logger.Log.Error("Error by save metrics to file", zap.Error(err))
	}

	appStorage.Close()

	logger.Log.Info("Terminate app")

	return nil
}

func startServerPprof() {
	http.ListenAndServe(flags.Args.AddrPprof, nil)
}
