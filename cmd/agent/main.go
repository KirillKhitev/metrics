// Агент для отправки метрик на сервер в формате JSON.
// Слепок метрик отправляется целиком.
// При ошибке отправки делаем 4 попытки. Между попытками - 2, 3 и 5 секунд соответственно.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/KirillKhitev/metrics/internal/client"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/KirillKhitev/metrics/internal/flags"
	"github.com/KirillKhitev/metrics/internal/metrics"
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
	client          client.Client
	runtimeStats    runtime.MemStats
	memStats        mem.SwapMemoryStat
	cpuStats        []float64
	pollCount       int64
	dataChan        chan []metrics.Metrics
	closeSenderChan chan struct{}
	wg              *sync.WaitGroup
}

// NewAgent конструктор главной структуры приложения Агента.
func NewAgent() (*agent, error) {
	client, err := newClient()
	if err != nil {
		return nil, err
	}

	a := &agent{
		client:          client,
		runtimeStats:    runtime.MemStats{},
		memStats:        mem.SwapMemoryStat{},
		dataChan:        make(chan []metrics.Metrics),
		closeSenderChan: make(chan struct{}),
		wg:              &sync.WaitGroup{},
	}

	return a, nil
}

func (a *agent) getMetricsFromRuntime() {
	ticker := time.Tick(time.Second * time.Duration(flags.ArgsClient.PollInterval))

	for {
		<-ticker
		runtime.ReadMemStats(&a.runtimeStats)
		a.pollCount++
	}
}

func (a *agent) getOtherMetrics() {
	ticker := time.Tick(time.Second * time.Duration(flags.ArgsClient.PollInterval))

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

func (a *agent) sender(ctx context.Context, idSender int) {
	for {
		select {
		case <-a.closeSenderChan:
			a.wg.Done()
			log.Printf("Stop Sender #%d", idSender)
			return
		default:
			select {
			case batchData := <-a.dataChan:
				var errSending error
				for i := 1; i <= AttemptCount; i++ {
					errSending = a.client.Send(ctx, batchData)
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

func (a *agent) workSenders() {
	ticker := time.Tick(time.Second * time.Duration(flags.ArgsClient.ReportInterval))

	for {
		<-ticker

		data := metrics.PrepareRuntimeMetrics(&a.runtimeStats)
		data = append(data, metrics.PrepareMemstatsMetrics(&a.memStats)...)
		data = append(data, metrics.PrepareCPUMetrics(a.cpuStats)...)
		data = append(data, metrics.PreparePollCounterMetric(a.pollCount))

		a.dataChan <- data
	}
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
	log.Println("Close data-channel")

	a.client.Close()

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
	flags.ArgsClient.ParseFlags()

	agent, err := NewAgent()
	if err != nil {
		log.Fatal(err)
	}

	agent.printBuildInfo()

	go agent.getMetricsFromRuntime()
	go agent.getOtherMetrics()

	ctx := context.Background()

	for w := 1; w <= flags.ArgsClient.RateLimit; w++ {
		agent.wg.Add(1)
		go agent.sender(ctx, w)
	}

	go agent.workSenders()

	if err := agent.catchTerminateSignal(); err != nil {
		log.Fatal(err)
	}
}

func newClient() (client.Client, error) {
	if flags.ArgsClient.GRPC {
		return client.NewGRPCClient()
	}

	return client.NewRestyClient()
}
