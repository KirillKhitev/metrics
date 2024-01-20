package handlers

import (
	"encoding/json"
	"github.com/KirillKhitev/metrics/internal/flags"
	"github.com/KirillKhitev/metrics/internal/logger"
	"github.com/KirillKhitev/metrics/internal/metrics"
	"go.uber.org/zap"
	"net/http"
)

type UpdatesHandler MyHandler

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

	if err := ch.Storage.UpdateCounters(r.Context(), counters); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := ch.Storage.UpdateGauges(r.Context(), gauges); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := make([]metrics.Metrics, 0)

	if err := ch.prepareResponseCounters(&response, r, counters); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := ch.prepareResponseGauges(&response, r, gauges); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if flags.Args.StoreInterval == 0 {
		if err := ch.Storage.TrySaveToFile(); err != nil {
			logger.Log.Error("Error by save metrics to file", zap.Error(err))
		}
	}

	str, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		logger.Log.Error("cannot encode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(str)
	w.WriteHeader(http.StatusOK)
}

func (ch *UpdatesHandler) prepareResponseCounters(response *[]metrics.Metrics, r *http.Request, counters []metrics.Metrics) error {
	var names []string

	for _, metrica := range counters {
		names = append(names, metrica.ID)
	}

	data, err := ch.Storage.GetCounters(r.Context(), names)
	if err != nil {
		return err
	}

	for name, value := range data {
		val := value
		metrica := metrics.Metrics{
			ID:    name,
			MType: "counter",
			Delta: &val,
		}
		*response = append(*response, metrica)
	}

	return nil
}

func (ch *UpdatesHandler) prepareResponseGauges(response *[]metrics.Metrics, r *http.Request, gauges []metrics.Metrics) error {
	var names []string

	for _, metrica := range gauges {
		names = append(names, metrica.ID)
	}

	data, err := ch.Storage.GetGauges(r.Context(), names)
	if err != nil {
		return err
	}

	for name, value := range data {
		val := value
		metrica := metrics.Metrics{
			ID:    name,
			MType: "gauges",
			Value: &val,
		}
		*response = append(*response, metrica)
	}

	return nil
}

func (ch *UpdatesHandler) getDataFromRequest(r *http.Request) ([]metrics.Metrics, []metrics.Metrics, error) {
	var request []metrics.Metrics

	counters := []metrics.Metrics{}
	gauges := []metrics.Metrics{}

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&request); err != nil {
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
