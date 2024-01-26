package handlers

import (
	"github.com/KirillKhitev/metrics/internal/logger"
	"go.uber.org/zap"
	"net/http"
)

type PingHandler MyHandler

func (ch *PingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err := ch.Storage.Ping(r.Context()); err != nil {
		logger.Log.Error("Cannot connect to db", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
