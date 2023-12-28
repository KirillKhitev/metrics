package handlers

import (
	"github.com/KirillKhitev/metrics/internal/storage"
	"net/http"
	"strconv"
)

type MyHandler struct {
	Storage storage.MemStorage
}

func updateCounter(ch *MyHandler, w http.ResponseWriter, n string, v string) bool {
	value, err := strconv.ParseInt(v, 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}

	err = ch.Storage.UpdateCounter(n, value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}

	w.WriteHeader(http.StatusOK)
	return true
}

func updateGauge(ch *MyHandler, w http.ResponseWriter, n string, v string) bool {
	value, err := strconv.ParseFloat(v, 64)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}

	err = ch.Storage.UpdateGauge(n, value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}

	w.WriteHeader(http.StatusOK)
	return true
}
