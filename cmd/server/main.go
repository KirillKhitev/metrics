// Сервер для приема метрик.
package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KirillKhitev/metrics/internal/mycrypto"

	"go.uber.org/zap"

	"github.com/KirillKhitev/metrics/internal/flags"
	"github.com/KirillKhitev/metrics/internal/gzip"
	"github.com/KirillKhitev/metrics/internal/logger"
	"github.com/KirillKhitev/metrics/internal/server"
	"github.com/KirillKhitev/metrics/internal/signature"
	"github.com/KirillKhitev/metrics/internal/storage"
)

// Флаги сборки
var (
	buildVersion string = "N/A" // Версия сборки
	buildDate    string = "N/A" // Дата сборки
	buildCommit  string = "N/A" // Комментарий сборки
)

func main() {
	printBuildInfo()
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

	srv := prepareHTTPServer(appStorage)

	go intervalSaveToFile(appStorage)
	go startServer(srv)
	go startServerPprof()

	return catchTerminateSignal(appStorage, srv)
}

// prepareHTTPServer создает http-сервер.
func prepareHTTPServer(appStorage storage.Repository) *http.Server {
	handler := mycrypto.Middleware(server.GetRouter(appStorage))
	handler = gzip.Middleware(handler)
	handler = signature.Middleware(handler)

	var srv = &http.Server{
		Addr:    flags.Args.AddrRun,
		Handler: handler,
	}

	return srv
}

// startServer запускает http-сервер.
func startServer(srv *http.Server) error {
	logger.Log.Info("Running server", zap.String("address", flags.Args.AddrRun))

	return srv.ListenAndServe()
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

// catchTerminateSignal ловит сигналы ОС для корректной остановки приложения.
func catchTerminateSignal(appStorage storage.Repository, srv *http.Server) error {
	terminateSignals := make(chan os.Signal, 1)

	signal.Notify(terminateSignals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-terminateSignals

	if err := shutdownHTTPServer(srv); err != nil {
		return err
	}

	if err := appStorage.TrySaveToFile(); err != nil {
		logger.Log.Error("Error by save metrics to file", zap.Error(err))
	}

	appStorage.Close()

	logger.Log.Info("Successful stop app server")

	return nil
}

// shutdownHTTPServer корректно останавливает http-сервер.
func shutdownHTTPServer(srv *http.Server) error {
	shutdownCtx, shutdownRelease := context.WithCancel(context.TODO())
	defer shutdownRelease()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("HTTP shutdown error: %w", err)
	}

	logger.Log.Info("Shutdown http-server")

	return nil
}

func startServerPprof() {
	http.ListenAndServe(flags.Args.AddrPprof, nil)
}

// printBuildInfo выводит в консоль информацию по сборке.
func printBuildInfo() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}
