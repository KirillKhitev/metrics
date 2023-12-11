package main

import (
	"fmt"
	"github.com/KirillKhitev/metrics/internal/metrics"
	"github.com/go-resty/resty/v2"
	"runtime"
	"sync"
	"time"
)

func main() {
	parseFlags()

	var PollCount int64
	var memStats runtime.MemStats

	mu := new(sync.Mutex)

	go func() {
		tickerUpdateMetrics := time.Tick(time.Second * time.Duration(flagPollInterval))

		for {
			<-tickerUpdateMetrics
			mu.Lock()
			runtime.ReadMemStats(&memStats)
			PollCount++
			mu.Unlock()
		}
	}()

	tickerSendMetrics := time.Tick(time.Second * time.Duration(flagReportInterval))

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
		_, err := sendUpdate(client, "counter", name, value)

		if err != nil {
			fmt.Println(err)
		}
	}

	for name, value := range metrics.PrepareGaugeForSend(m) {
		_, err := sendUpdate(client, "gauge", name, value)

		if err != nil {
			fmt.Println(err)
		}
	}
}

func sendUpdate(client *resty.Client, t, name, value string) (*resty.Response, error) {
	url := fmt.Sprintf("http://%s/update/%s/%s/%s", flagAddrRun, t, name, value)
	resp, err := client.R().Post(url)

	return resp, err
}
