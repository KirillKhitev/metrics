package handlers

import (
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/KirillKhitev/metrics/internal/flags"
	"github.com/KirillKhitev/metrics/internal/logger"
)

// UpdateHandler отвечает за обработку POST-запроса /update/{type}/{name}.
type UpdateHandler MyHandler

/*
ServeHTTP служит для добавления/обновления отдельной метрики.

Коды ответа:

• 200 - успешный запрос.

• 400 - неуказаны тип или значение метрики, или тип метрики неподдерживается.

• 404 - неуказано название метрики.

• 405 - метод запроса отличен от POST.
*/
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
		res = updateCounter(r.Context(), &mh, w, nameMetric, valueMetric)
	case "gauge":
		res = updateGauge(r.Context(), &mh, w, nameMetric, valueMetric)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}

	if !res {
		return
	}

	if flags.Args.StoreInterval == 0 {
		if err := ch.Storage.TrySaveToFile(); err != nil {
			logger.Log.Error("Error by save metrics to file", zap.Error(err))
		}
	}

	w.WriteHeader(http.StatusOK)
}
