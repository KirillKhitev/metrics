package handlers

import (
	"fmt"
	"net/http"
)

type ListHandler MyHandler

func (ch *ListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	result := "<b>Counter:</b><br/>"

	for name, value := range ch.Storage.GetCounterList() {
		result += fmt.Sprintf("<p>%s: %d</p>", name, value)
	}

	result += "<br/><br/><b>Gauge:</b><br/>"

	for name, value := range ch.Storage.GetGaugeList() {
		result += fmt.Sprintf("<p>%s: %g</p>", name, value)
	}

	w.Write([]byte(result))
	w.WriteHeader(http.StatusOK)
}
