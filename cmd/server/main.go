package main

import (
	"github.com/KirillKhitev/metrics/internal/dump"
	"github.com/KirillKhitev/metrics/internal/flags"
	"github.com/KirillKhitev/metrics/internal/gzip"
	"github.com/KirillKhitev/metrics/internal/logger"
	"github.com/KirillKhitev/metrics/internal/server"
	"github.com/KirillKhitev/metrics/internal/storage"
	"go.uber.org/zap"
	"net/http"
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

	appStorage := storage.MemStorage{}
	if err := appStorage.Init(); err != nil {
		return err
	}

	if flags.Args.StoreInterval > 0 {
		go func(appStorage storage.MemStorage) {
			ticker := time.Tick(time.Second * time.Duration(flags.Args.StoreInterval))

			for {
				<-ticker
				dump.SaveStorageToFile(flags.Args.FileStoragePath, appStorage)
			}
		}(appStorage)
	}

	defer dump.SaveStorageToFile(flags.Args.FileStoragePath, appStorage)

	logger.Log.Info("Running server", zap.String("address", flags.Args.AddrRun))

	return http.ListenAndServe(flags.Args.AddrRun, logger.RequestLogger(gzip.Middleware(server.GetRouter(appStorage))))
}
