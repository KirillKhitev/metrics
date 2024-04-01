package handlers

import (
	"fmt"
	"net/http"
	"strings"
)

// GetHandler отвечает за обработку GET-запроса /value/{type}/{name}.
type GetHandler MyHandler

/*
ServeHTTP возвращает значение метрики по ее названию и типу.

Коды ответа:

• 200 - успешный запрос.

• 400 - неуказан тип метрики, либо этот тип не поддерживается.

• 404 - неуказанно название метрики, либо по данной метрике нет данных.

• 405 - метод запроса отличен от GET.
*/
func (ch *GetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
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

	nameMetric := pathParts[3]
	res := ""

	switch typeMetric {
	case "counter":
		v, err := ch.Storage.GetCounter(r.Context(), nameMetric)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		res = fmt.Sprintf("%d", v)
	case "gauge":
		v, err := ch.Storage.GetGauge(r.Context(), nameMetric)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		res = fmt.Sprintf("%g", v)
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write([]byte(res))
	w.WriteHeader(http.StatusOK)
}
