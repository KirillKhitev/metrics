// Агент для отправки метрик на сервер в формате JSON.
// Слепок метрик отправляется целиком.
// При ошибке отправки делаем 4 попытки. Между попытками - 2, 3 и 5 секунд соответственно.
package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/KirillKhitev/metrics/internal/mycrypto"

	"github.com/go-resty/resty/v2"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/KirillKhitev/metrics/internal/metrics"
	"github.com/KirillKhitev/metrics/internal/signature"
)

// AttemptCount определяет количество попыток отправки данных на сервер
const AttemptCount = 4

// Флаги сборки
var (
	buildVersion string = "N/A" // Версия сборки
	buildDate    string = "N/A" // Дата сборки
	buildCommit  string = "N/A" // Комментарий сборки
)

type agent struct {
	client          *resty.Client
	runtimeStats    runtime.MemStats
	memStats        mem.SwapMemoryStat
	cpuStats        []float64
	pollCount       int64
	dataChan        chan []metrics.Metrics
	closeSenderChan chan struct{}
	wg              *sync.WaitGroup
}

// NewAgent конструктор главной структуры приложения Агента.
func NewAgent() *agent {
	return &agent{
		client:          resty.New(),
		runtimeStats:    runtime.MemStats{},
		memStats:        mem.SwapMemoryStat{},
		dataChan:        make(chan []metrics.Metrics),
		closeSenderChan: make(chan struct{}),
		wg:              &sync.WaitGroup{},
	}
}

// Compress сжимает данные перед отправкой на сервер.
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

func (a *agent) getMetricsFromRuntime() {
	ticker := time.Tick(time.Second * time.Duration(flags.PollInterval))

	for {
		<-ticker
		runtime.ReadMemStats(&a.runtimeStats)
		a.pollCount++
	}
}

func (a *agent) getOtherMetrics() {
	ticker := time.Tick(time.Second * time.Duration(flags.PollInterval))

	for {
		<-ticker
		memStats, err := mem.SwapMemory()
		if err != nil {
			log.Printf("Failure get metrics SwapMemory: %s", err)
		}
		a.memStats = *memStats

		cpuStats, err := cpu.Percent(0, true)
		if err != nil {
			log.Printf("Failure get metrics cpuStats: %s", err)
			return
		}

		a.cpuStats = cpuStats
	}
}

func (a *agent) sender(idSender int) {
	for {
		select {
		case <-a.closeSenderChan:
			a.wg.Done()
			log.Printf("Stop Sender #%d", idSender)
			return
		default:
			select {
			case batchData := <-a.dataChan:
				data, err := a.prepareDataForSend(batchData, idSender)
				if err != nil {
					log.Printf("sender %d, failure prepare data for send, err: %s", idSender, err)
					continue
				}

				var errSending error
				for i := 1; i <= AttemptCount; i++ {
					_, errSending = a.send(data)
					if errSending != nil {
						log.Printf("sender %d, attempt%d send metrics, err: %s", idSender, i, errSending)

						if i < AttemptCount {
							time.Sleep(time.Duration(2*i-1) * time.Second)
						}

						continue
					}

					break
				}

				if errSending != nil {
					log.Printf("sender %d, failure send metrics", idSender)
				}
			default:
			}
		}
	}
}

// prepareDataForSend готовит данные для отправки на сервер.
func (a *agent) prepareDataForSend(batchData []metrics.Metrics, idSender int) ([]byte, error) {
	data, err := json.Marshal(batchData)
	if err != nil {
		return data, fmt.Errorf("error by json-encode metrics: %w", err)
	}

	dataEncrypted, err := mycrypto.Encrypting(data, flags.CryptoKey)
	if err != nil {
		return data, fmt.Errorf("error by encrypting data: %s, err: %w", data, err)
	}

	dataCompress, err := a.Compress(dataEncrypted)
	if err != nil {
		return data, fmt.Errorf("error by compress data: %s, err: %w", data, err)
	}

	return dataCompress, nil
}

func (a *agent) workSenders() {
	ticker := time.Tick(time.Second * time.Duration(flags.ReportInterval))

	for {
		<-ticker

		data := metrics.PrepareRuntimeMetrics(&a.runtimeStats)
		data = append(data, metrics.PrepareMemstatsMetrics(&a.memStats)...)
		data = append(data, metrics.PrepareCPUMetrics(a.cpuStats)...)
		data = append(data, metrics.PreparePollCounterMetric(a.pollCount))

		a.dataChan <- data
	}
}

func (a *agent) send(data []byte) (*resty.Response, error) {
	url := fmt.Sprintf("http://%s/updates/", flags.AddrRun)

	request := a.client.NewRequest().
		SetBody(data).
		SetHeader("Content-Encoding", "gzip")

	if flags.Key != "" {
		hashSum := signature.GetHash(data, flags.Key)
		request.SetHeader("HashSHA256", hashSum)
	}

	resp, err := request.Post(url)

	return resp, err
}

// printBuildInfo выводит в консоль информацию по сборке.
func (a *agent) printBuildInfo() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}

// catchTerminateSignal ловит сигналы ОС для корректной остановки агента.
func (a *agent) catchTerminateSignal() error {
	terminateSignals := make(chan os.Signal, 1)

	signal.Notify(terminateSignals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-terminateSignals

	if err := a.Close(); err != nil {
		return err
	}

	return nil
}

// Close отвечает за корректную остановку агента.
func (a *agent) Close() error {
	a.stopSenders()

	close(a.dataChan)

	log.Println("Successful stop agent")

	return nil
}

// stopSenders останавливает воркеры.
func (a *agent) stopSenders() {
	log.Println("Waiting closing all senders")

	close(a.closeSenderChan)
	a.wg.Wait()

	log.Println("All senders are stopped!")
}

func main() {
	agent := NewAgent()
	agent.printBuildInfo()

	flags.ParseFlags()

	go agent.getMetricsFromRuntime()
	go agent.getOtherMetrics()

	for w := 1; w <= flags.RateLimit; w++ {
		agent.wg.Add(1)
		go agent.sender(w)
	}

	go agent.workSenders()

	if err := agent.catchTerminateSignal(); err != nil {
		log.Fatal(err)
	}
}
