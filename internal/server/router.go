package server

import (
	"github.com/KirillKhitev/metrics/internal/handlers"
	"github.com/KirillKhitev/metrics/internal/storage"
	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	Update handlers.UpdateHandler
	List   handlers.ListHandler
	Get    handlers.GetHandler
}

func GetRouter() chi.Router {
	r := chi.NewRouter()

	appStorage := storage.MemStorage{}
	if err := appStorage.Init(); err != nil {
		panic(err)
	}

	var myHandlers = Handlers{
		Update: handlers.UpdateHandler{
			Storage: appStorage,
		},
		List: handlers.ListHandler{
			Storage: appStorage,
		},
		Get: handlers.GetHandler{
			Storage: appStorage,
		},
	}

	r.Route("/", func(r chi.Router) {
		r.Handle("/", &myHandlers.List)
		r.Route("/update", func(r chi.Router) {
			r.Post("/", handlers.BadRequestHandle)
			r.Route("/{typeMetric}", func(r chi.Router) {
				r.Post("/", handlers.NotFoundHandle)
				r.Route("/{nameMetric}", func(r chi.Router) {
					r.Post("/", handlers.BadRequestHandle)
					r.Handle("/{valueMetric}", &myHandlers.Update)
				})
			})
		})
		r.Route("/value", func(r chi.Router) {
			r.Get("/", handlers.BadRequestHandle)
			r.Route("/{typeMetric}", func(r chi.Router) {
				r.Get("/", handlers.NotFoundHandle)
				r.Handle("/{nameMetric}", &myHandlers.Get)
			})
		})
	})

	return r
}
