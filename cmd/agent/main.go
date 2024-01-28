package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/KirillKhitev/metrics/internal/metrics"
	"github.com/KirillKhitev/metrics/internal/signature"
	"github.com/go-resty/resty/v2"
	"log"
	"runtime"
	"sync"
	"time"
)

const AttemptCount = 4

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

	for {
		<-ticker
		a.Lock()

		data, err := a.prepareDataForSend()
		if err != nil {
			log.Printf("Failure prepare data for send: %s", err)
			a.Unlock()
			continue
		}

		for i := 1; i <= AttemptCount; i++ {
			_, err = a.sendUpdate(data)
			if err != nil {
				log.Printf("Attempt%d send metrics, err: %v", i, err)

				if i < AttemptCount {
					time.Sleep(time.Duration(2*i-1) * time.Second)
				}

				continue
			}

			break
		}

		if err != nil {
			log.Println("Failure send metrics")
		}

		a.Unlock()
	}
}

func (a *agent) prepareDataForSend() ([]byte, error) {
	data := make([]metrics.Metrics, 0)

	for name, value := range metrics.PrepareCounterForSend(a.pollCount) {
		valueMetric := value
		metrica := metrics.Metrics{
			ID:    name,
			MType: "counter",
			Delta: &valueMetric,
		}

		data = append(data, metrica)
	}

	for name, value := range metrics.PrepareGaugeForSend(&a.data) {
		valueMetric := value
		metrica := metrics.Metrics{
			ID:    name,
			MType: "gauge",
			Value: &valueMetric,
		}

		data = append(data, metrica)
	}

	result, err := json.Marshal(data)
	if err != nil {
		err = fmt.Errorf("error by encode metric: %v, error: %w", data, err)
		return result, err
	}

	return result, nil
}

func (a *agent) sendUpdate(data []byte) (*resty.Response, error) {
	url := fmt.Sprintf("http://%s/updates/", flags.AddrRun)

	dataCompress, err := a.Compress(data)
	if err != nil {
		return nil, err
	}

	request := a.client.NewRequest().
		SetBody(dataCompress).
		SetHeader("Content-Encoding", "gzip")

	if flags.Key != "" {
		hashSum := signature.GetHash(dataCompress, flags.Key)
		request.SetHeader("HashSHA256", hashSum)
	}

	resp, err := request.Post(url)

	return resp, err
}

func main() {
	agent := agent{}
	agent.client = resty.New()

	flags.ParseFlags()

	go agent.getMetrics()

	agent.sendMetrics()
}
