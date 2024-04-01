package handlers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/KirillKhitev/metrics/internal/logger"
	"github.com/KirillKhitev/metrics/internal/metrics"
)

// GetJSONHandler отвечает за обработку POST-запроса /value.
type GetJSONHandler MyHandler

/*
ServeHTTP возвращает json-представление объекта Метрики.

Коды ответа:

• 200 - успешный запрос.

• 400 - метод запроса отличен от POST, или не передан тип метрики или он не поддерживается.

• 404 - неуказанно название метрики, либо по данной метрике нет данных.

• 500 - при возникновении ошибки.
*/
func (ch *GetJSONHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

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

	switch request.MType {
	case "counter":
		v, err := ch.Storage.GetCounter(r.Context(), request.ID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		request.Delta = &v
	case "gauge":
		v, err := ch.Storage.GetGauge(r.Context(), request.ID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		request.Value = &v
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	str, err := json.MarshalIndent(request, "", "    ")
	if err != nil {
		logger.Log.Error("cannot encode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(str)
	w.WriteHeader(http.StatusOK)
}
