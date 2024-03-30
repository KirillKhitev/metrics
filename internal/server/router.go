package server

import (
	"github.com/go-chi/chi/v5"

	"github.com/KirillKhitev/metrics/internal/handlers"
	"github.com/KirillKhitev/metrics/internal/storage"
)

type Handlers struct {
	Update     handlers.UpdateHandler
	UpdateJSON handlers.UpdateJSONHandler
	Updates    handlers.UpdatesHandler
	List       handlers.ListHandler
	Get        handlers.GetHandler
	GetJSON    handlers.GetJSONHandler
	Ping       handlers.PingHandler
}

func GetRouter(appStorage storage.Repository) chi.Router {
	r := chi.NewRouter()

	var myHandlers = Handlers{
		Update: handlers.UpdateHandler{
			Storage: appStorage,
		},
		UpdateJSON: handlers.UpdateJSONHandler{
			Storage: appStorage,
		},
		Updates: handlers.UpdatesHandler{
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
		Ping: handlers.PingHandler{
			Storage: appStorage,
		},
	}

	r.Route("/", func(r chi.Router) {
		r.Handle("/", &myHandlers.List)
		r.Route("/updates", func(r chi.Router) {
			r.Handle("/", &myHandlers.Updates)
		})
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
		r.Handle("/ping", &myHandlers.Ping)
	})

	return r
}
