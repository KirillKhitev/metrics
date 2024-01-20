package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/KirillKhitev/metrics/internal/metrics"
	"github.com/go-resty/resty/v2"
	"log"
	"runtime"
	"sync"
	"time"
)

type agent struct {
	sync.Mutex
	client    *resty.Client
	data      runtime.MemStats
	pollCount int64
}

func (a *agent) Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer

	w := gzip.NewWriter(&b)

	_, err := w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed write data to compress temporary buffer: %v", err)
	}

	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed compress data: %v", err)
	}

	return b.Bytes(), nil
}

func (a *agent) getMetrics() {
	ticker := time.Tick(time.Second * time.Duration(flags.PollInterval))

	for {
		<-ticker
		a.Lock()
		runtime.ReadMemStats(&a.data)
		a.pollCount++
		a.Unlock()
	}
}

func (a *agent) sendMetrics() {
	ticker := time.Tick(time.Second * time.Duration(flags.ReportInterval))

	send := func(body []metrics.Metrics) {
		str, err := json.Marshal(body)
		if err != nil {
			log.Printf("error by encode metric: %v, error: %s", body, err)
			return
		}

		_, err = a.sendUpdate(str)

		if err != nil {
			log.Println(err)
			return
		}
	}

	for {
		<-ticker
		a.Lock()

		data := make([]metrics.Metrics, 0)

		for name, value := range metrics.PrepareCounterForSend(a.pollCount) {
			metrica := metrics.Metrics{
				ID:    name,
				MType: "counter",
				Delta: &value,
			}

			data = append(data, metrica)
		}

		for name, value := range metrics.PrepareGaugeForSend(&a.data) {
			metrica := metrics.Metrics{
				ID:    name,
				MType: "gauge",
				Value: &value,
			}

			data = append(data, metrica)
		}

		send(data)

		a.Unlock()
	}
}

func (a *agent) sendUpdate(data []byte) (*resty.Response, error) {
	url := fmt.Sprintf("http://%s/updates/", flags.AddrRun)

	dataCompress, err := a.Compress(data)
	if err != nil {
		return nil, err
	}

	resp, err := a.client.R().
		SetBody(dataCompress).
		SetHeader("Content-Encoding", "gzip").
		Post(url)

	return resp, err
}

func main() {
	agent := agent{}
	agent.client = resty.New()

	flags.ParseFlags()

	go agent.getMetrics()

	agent.sendMetrics()
}
