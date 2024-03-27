package handlers

import (
	"bytes"
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

	buf := bytes.Buffer{}
	buf.WriteString("<b>Counter:</b><br/>")

	for name, value := range ch.Storage.GetCounterList(r.Context()) {
		buf.WriteString(fmt.Sprintf("<p>%s: %d</p>", name, value))
	}

	buf.WriteString("<br/><br/><b>Gauge:</b><br/>")

	for name, value := range ch.Storage.GetGaugeList(r.Context()) {
		buf.WriteString(fmt.Sprintf("<p>%s: %g</p>", name, value))
	}

	w.Write(buf.Bytes())
	w.WriteHeader(http.StatusOK)
}
