package handlers

import (
	"context"
	"github.com/KirillKhitev/metrics/internal/storage"
	"net/http"
	"strconv"
)

type MyHandler struct {
	Storage storage.Repository
}

func updateCounter(ctx context.Context, ch *MyHandler, w http.ResponseWriter, name string, valStr string) bool {
	value, err := strconv.ParseInt(valStr, 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}

	err = ch.Storage.UpdateCounter(ctx, name, value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}

	w.WriteHeader(http.StatusOK)
	return true
}

func updateGauge(ctx context.Context, ch *MyHandler, w http.ResponseWriter, name string, valStr string) bool {
	value, err := strconv.ParseFloat(valStr, 64)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}

	err = ch.Storage.UpdateGauge(ctx, name, value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}

	w.WriteHeader(http.StatusOK)
	return true
}
