package main

import (
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

	send := func(body metrics.Metrics) {
		str, err := json.Marshal(body)
		if err != nil {
			log.Printf("error by encode metric: %v, error: %s", body, err)
			return
		}

		_, err = a.sendUpdate(string(str))

		if err != nil {
			log.Println(err)
			return
		}
	}

	for {
		<-ticker
		a.Lock()

		for name, value := range metrics.PrepareCounterForSend(a.pollCount) {
			metrica := metrics.Metrics{
				ID:    name,
				MType: "counter",
				Delta: &value,
			}

			send(metrica)
		}

		for name, value := range metrics.PrepareGaugeForSend(&a.data) {
			metrica := metrics.Metrics{
				ID:    name,
				MType: "gauge",
				Value: &value,
			}

			send(metrica)
		}

		a.Unlock()
	}
}

func (a *agent) sendUpdate(data string) (*resty.Response, error) {
	url := fmt.Sprintf("http://%s/update/", flags.AddrRun)

	resp, err := a.client.R().
		SetBody(data).
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
