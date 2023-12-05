package server

import (
	"github.com/KirillKhitev/metrics/internal/handlers"
	"github.com/KirillKhitev/metrics/internal/storage"
	"net/http"
)

func Run() error {
	appStorage := storage.MemStorage{}
	if err := appStorage.Init(); err != nil {
		return err
	}

	mux := http.NewServeMux()
	updateHandler := &handlers.UpdateHandler{
		Storage: appStorage,
	}

	mux.Handle(`/update/`, updateHandler)
	mux.HandleFunc(`/`, handlers.AllPageHandler)

	return http.ListenAndServe(`:8080`, mux)
}
