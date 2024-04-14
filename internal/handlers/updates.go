package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/KirillKhitev/metrics/internal/flags"
	"github.com/KirillKhitev/metrics/internal/logger"
	"github.com/KirillKhitev/metrics/internal/metrics"
)

// UpdatesHandler отвечает за обработку POST-запроса /updates.
type UpdatesHandler MyHandler

/*
ServeHTTP служит для добавления/обновления списка метрик.

Коды ответа:

• 200 - успешный запрос.

• 400 - неверный запрос.

• 405 - метод запроса отличен от POST.

• 500 - при возникновении внутренней ошибки.
*/
func (ch *UpdatesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	counters, gauges, err := ch.getDataFromRequest(r)
	if err != nil {
		logger.Log.Error("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if errUpdate := ch.Storage.UpdateCounters(r.Context(), counters); errUpdate != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if errUpdate := ch.Storage.UpdateGauges(r.Context(), gauges); errUpdate != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := make([]metrics.Metrics, 0, len(counters)+len(gauges))

	if response, err = ch.prepareResponseCounters(response, r, counters); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if response, err = ch.prepareResponseGauges(response, r, gauges); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if flags.Args.StoreInterval == 0 {
		if errSaveFile := ch.Storage.TrySaveToFile(); errSaveFile != nil {
			logger.Log.Error("Error by save metrics to file", zap.Error(errSaveFile))
		}
	}

	str, errJSON := json.MarshalIndent(response, "", "    ")
	if errJSON != nil {
		logger.Log.Error("cannot encode request JSON body", zap.Error(errJSON))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(str)
	w.WriteHeader(http.StatusOK)
}

// prepareResponseCounters готовит список метрик с типом counter из хранилища.
func (ch *UpdatesHandler) prepareResponseCounters(response []metrics.Metrics, r *http.Request, counters []metrics.Metrics) ([]metrics.Metrics, error) {
	var names []string

	for _, metrica := range counters {
		names = append(names, metrica.ID)
	}

	data, err := ch.Storage.GetCounters(r.Context(), names)
	if err != nil {
		return response, err
	}

	for name, value := range data {
		val := value
		metrica := metrics.Metrics{
			ID:    name,
			MType: "counter",
			Delta: &val,
		}
		response = append(response, metrica)
	}

	return response, nil
}

// prepareResponseGauges готовит список метрик с типом gauge из Хранилища.
func (ch *UpdatesHandler) prepareResponseGauges(response []metrics.Metrics, r *http.Request, gauges []metrics.Metrics) ([]metrics.Metrics, error) {
	var names []string

	for _, metrica := range gauges {
		names = append(names, metrica.ID)
	}

	data, err := ch.Storage.GetGauges(r.Context(), names)
	if err != nil {
		return response, err
	}

	for name, value := range data {
		val := value
		metrica := metrics.Metrics{
			ID:    name,
			MType: "gauges",
			Value: &val,
		}
		response = append(response, metrica)
	}

	return response, nil
}

// getDataFromRequest формирует список метрик из Запроса.
func (ch *UpdatesHandler) getDataFromRequest(r *http.Request) ([]metrics.Metrics, []metrics.Metrics, error) {
	var request []metrics.Metrics

	counters := []metrics.Metrics{}
	gauges := []metrics.Metrics{}

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&request); err != nil && err != io.EOF {
		return counters, gauges, err
	}

	for _, metrica := range request {
		if metrica.MType == "" || metrica.ID == "" || (metrica.Value == nil && metrica.Delta == nil) {
			continue
		}

		switch metrica.MType {
		case "counter":
			counters = append(counters, metrica)
		case "gauge":
			gauges = append(gauges, metrica)
		}
	}

	return counters, gauges, nil
}
