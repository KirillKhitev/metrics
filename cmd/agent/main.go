package main

import (
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

	for {
		<-ticker
		a.Lock()

		for name, value := range metrics.PrepareCounterForSend(a.pollCount) {
			_, err := a.sendUpdate("counter", name, value)

			if err != nil {
				log.Println(err)
			}
		}

		for name, value := range metrics.PrepareGaugeForSend(&a.data) {
			_, err := a.sendUpdate("gauge", name, value)

			if err != nil {
				log.Println(err)
			}
		}

		a.Unlock()
	}
}

func (a *agent) sendUpdate(typeMetric, name, value string) (*resty.Response, error) {
	url := fmt.Sprintf("http://%s/update/%s/%s/%s", flags.AddrRun, typeMetric, name, value)
	resp, err := a.client.R().Post(url)

	return resp, err
}

func main() {
	agent := agent{}
	agent.client = resty.New()

	flags.ParseFlags()

	go agent.getMetrics()

	agent.sendMetrics()
}
