package handlers

import (
	"context"
	"github.com/KirillKhitev/metrics/internal/logger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type PingHandler MyHandler

func (ch *PingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := ch.Storage.DB.PingContext(ctx); err != nil {
		logger.Log.Error("Cannot connect to db", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
