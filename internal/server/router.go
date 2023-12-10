package server

import (
	"github.com/KirillKhitev/metrics/internal/handlers"
	"github.com/KirillKhitev/metrics/internal/storage"
	"github.com/go-chi/chi/v5"
)

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
			r.Post("/", handlers.BadRequestHandle)
			r.Route("/{typeMetric}", func(r chi.Router) {
				r.Post("/", handlers.NotFoundHandle)
				r.Route("/{nameMetric}", func(r chi.Router) {
					r.Post("/", handlers.BadRequestHandle)
					r.Handle("/{valueMetric}", updateHandler)
				})
			})
		})
		r.Route("/value", func(r chi.Router) {
			r.Get("/", handlers.BadRequestHandle)
			r.Route("/{typeMetric}", func(r chi.Router) {
				r.Get("/", handlers.NotFoundHandle)
				r.Handle("/{nameMetric}", valueHandler)
			})
		})
	})

	return r
}
