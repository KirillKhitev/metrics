package server

import (
	"github.com/KirillKhitev/metrics/internal/handlers"
	"github.com/KirillKhitev/metrics/internal/storage"
	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	Update     handlers.UpdateHandler
	UpdateJSON handlers.UpdateJSONHandler
	List       handlers.ListHandler
	Get        handlers.GetHandler
	GetJSON    handlers.GetJSONHandler
}

func GetRouter(appStorage storage.MemStorage) chi.Router {
	r := chi.NewRouter()

	var myHandlers = Handlers{
		Update: handlers.UpdateHandler{
			Storage: appStorage,
		},
		UpdateJSON: handlers.UpdateJSONHandler{
			Storage: appStorage,
		},
		List: handlers.ListHandler{
			Storage: appStorage,
		},
		Get: handlers.GetHandler{
			Storage: appStorage,
		},
		GetJSON: handlers.GetJSONHandler{
			Storage: appStorage,
		},
	}

	r.Route("/", func(r chi.Router) {
		r.Handle("/", &myHandlers.List)
		r.Route("/update", func(r chi.Router) {
			r.Handle("/", &myHandlers.UpdateJSON)
			r.Route("/{typeMetric}", func(r chi.Router) {
				r.Post("/", handlers.NotFoundHandle)
				r.Route("/{nameMetric}", func(r chi.Router) {
					r.Post("/", handlers.BadRequestHandle)
					r.Handle("/{valueMetric}", &myHandlers.Update)
				})
			})
		})
		r.Route("/value", func(r chi.Router) {
			r.Handle("/", &myHandlers.GetJSON)
			r.Route("/{typeMetric}", func(r chi.Router) {
				r.Get("/", handlers.NotFoundHandle)
				r.Handle("/{nameMetric}", &myHandlers.Get)
			})
		})
	})

	return r
}
