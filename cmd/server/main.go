package main

import (
	"github.com/KirillKhitev/metrics/internal/config"
	"github.com/KirillKhitev/metrics/internal/server"
	"log"
	"net/http"
)

func main() {
	log.Fatal(http.ListenAndServe(config.ServerHost+config.ServerPort, server.GetRouter()))
}

/*
func GetRouter() chi.Router {
    r := chi.NewRouter()

    appStorage := storage.MemStorage{}
    if err := appStorage.Init(); err != nil {
        panic(err)
    }

    updateHandler := &handlers.UpdateHandler{
        Storage: appStorage,
    }

    listHandler := &handlers.ListHandler{
        Storage: appStorage,
    }

    valueHandler := &handlers.ValueHandler{
        Storage: appStorage,
    }

    r.Route("/", func(r chi.Router) {
        r.Handle("/", listHandler)
        r.Route("/update", func(r chi.Router) {
            r.Post("/", BadRequestHandle)
            r.Route("/{typeMetric}", func(r chi.Router) {
                r.Post("/", NotFoundHandle)
                r.Route("/{nameMetric}", func(r chi.Router) {
                    r.Post("/", BadRequestHandle)
                    r.Handle("/{valueMetric}", updateHandler)
                })
            })
        })
        r.Route("/value", func(r chi.Router) {
            r.Get("/", BadRequestHandle)
            r.Route("/{typeMetric}", func(r chi.Router) {
                r.Get("/", NotFoundHandle)
                r.Handle("/{nameMetric}", valueHandler)
            })
        })
    })

    return r
}

func BadRequestHandle(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusBadRequest)
}

func NotFoundHandle(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusNotFound)
}
*/
