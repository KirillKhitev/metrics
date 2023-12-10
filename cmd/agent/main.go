package main

import (
	"fmt"
	"github.com/KirillKhitev/metrics/internal/agent"
	"github.com/KirillKhitev/metrics/internal/config"
	"github.com/KirillKhitev/metrics/internal/metrics"
	"github.com/go-resty/resty/v2"
	"runtime"
	"sync"
	"time"
)

func main() {
	var PollCount int64
	var memStats runtime.MemStats

	mu := new(sync.Mutex)

	go func() {
		tickerUpdateMetrics := time.Tick(time.Second * config.PollInterval)

		for {
			<-tickerUpdateMetrics
			mu.Lock()
			runtime.ReadMemStats(&memStats)
			PollCount++
			mu.Unlock()
		}
	}()

	tickerSendMetrics := time.Tick(time.Second * config.ReportInterval)

	for {
		<-tickerSendMetrics
		mu.Lock()

		sendDataToServer(&memStats, PollCount)

		mu.Unlock()
	}
}

func sendDataToServer(m *runtime.MemStats, PollCount int64) {
	client := resty.New()

	for name, value := range metrics.PrepareCounterForSend(PollCount) {
		_, err := agent.SendUpdate(client, "counter", name, value)

		if err != nil {
			fmt.Println(err)
		}
	}

	for name, value := range metrics.PrepareGaugeForSend(m) {
		_, err := agent.SendUpdate(client, "gauge", name, value)

		if err != nil {
			fmt.Println(err)
		}
	}
}
