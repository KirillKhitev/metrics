package handlers

import (
	"github.com/KirillKhitev/metrics/internal/storage"
	"net/http"
	"strconv"
	"strings"
)

type UpdateHandler struct {
	Storage storage.MemStorage
}

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

	switch typeMetric {
	case "counter":
		updateCounter(ch, w, nameMetric, valueMetric)
		break
	case "gauge":
		updateGauge(ch, w, nameMetric, valueMetric)
		break
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func AllPageHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func updateCounter(ch *UpdateHandler, w http.ResponseWriter, n string, v string) {
	value, err := strconv.ParseInt(v, 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = ch.Storage.UpdateCounter(n, value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func updateGauge(ch *UpdateHandler, w http.ResponseWriter, n string, v string) {
	value, err := strconv.ParseFloat(v, 64)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = ch.Storage.UpdateGauge(n, value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
