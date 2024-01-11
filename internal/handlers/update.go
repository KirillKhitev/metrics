package handlers

import (
	"github.com/KirillKhitev/metrics/internal/flags"
	"net/http"
	"strings"
)

type UpdateHandler MyHandler

func (ch *UpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	if len(pathParts) < 5 || pathParts[4] == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	nameMetric := pathParts[3]
	valueMetric := pathParts[4]

	mh := MyHandler(*ch)
	res := false

	switch typeMetric {
	case "counter":
		res = updateCounter(&mh, w, nameMetric, valueMetric)
	case "gauge":
		res = updateGauge(&mh, w, nameMetric, valueMetric)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}

	if !res {
		return
	}

	if flags.Args.StoreInterval == 0 {
		ch.Storage.SaveToFile()
	}

	w.WriteHeader(http.StatusOK)
}
