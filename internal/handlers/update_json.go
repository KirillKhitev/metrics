package handlers

import (
	"encoding/json"
	"github.com/KirillKhitev/metrics/internal/logger"
	"github.com/KirillKhitev/metrics/internal/metrics"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type UpdateJsonHandler MyHandler

func (ch *UpdateJsonHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var request metrics.Metrics

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&request); err != nil {
		logger.Log.Error("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if request.MType == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if request.ID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if request.Value == nil && request.Delta == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mh := MyHandler(*ch)

	w.Header().Set("Content-Type", "application/json")
	res := false

	switch request.MType {
	case "counter":
		res = updateCounter(&mh, w, request.ID, strconv.FormatInt(*request.Delta, 10))
		*request.Delta, _ = ch.Storage.GetCounter(request.ID)

	case "gauge":
		res = updateGauge(&mh, w, request.ID, strconv.FormatFloat(*request.Value, 'g', -1, 64))
		*request.Value, _ = ch.Storage.GetGauge(request.ID)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}

	if !res {
		return
	}

	str, err := json.MarshalIndent(request, "", "    ")
	if err != nil {
		logger.Log.Error("cannot encode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(str)
}
