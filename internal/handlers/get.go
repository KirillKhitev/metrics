package handlers

import (
	"fmt"
	"github.com/KirillKhitev/metrics/internal/storage"
	"net/http"
	"strings"
)

type GetHandler struct {
	Storage storage.MemStorage
}

func (ch *GetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	pathParts := strings.Split(r.RequestURI, "/")

	typeMetric := pathParts[2]

	if typeMetric == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(pathParts) < 4 || pathParts[3] == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	nameMetric := pathParts[3]
	res := ""

	switch typeMetric {
	case "counter":
		v, err := ch.Storage.GetCounter(nameMetric)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		res = fmt.Sprintf("%d", v)
	case "gauge":
		v, err := ch.Storage.GetGauge(nameMetric)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		res = fmt.Sprintf("%g", v)
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write([]byte(res))
}
