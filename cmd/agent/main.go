package main

import (
	"fmt"
	"github.com/KirillKhitev/metrics/internal/agent"
	"github.com/KirillKhitev/metrics/internal/config"
	"github.com/KirillKhitev/metrics/internal/metrics"
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

		counter := metrics.PrepareCounterForSend(PollCount)
		gauge := metrics.PrepareGaugeForSend(&memStats)

		for name, value := range counter {
			_, _ = agent.SendUpdate("counter", name, fmt.Sprintf("%d", value))
		}

		for name, value := range gauge {
			_, _ = agent.SendUpdate("gauge", name, fmt.Sprintf("%f", value))
		}

		mu.Unlock()
	}
}
