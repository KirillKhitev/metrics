// Модель объекта Метрика
package metrics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"reflect"
	"runtime"

	"github.com/shirou/gopsutil/v3/mem"
)

// Список runtime-метрик
var runtimeMetricsNames = []string{
	"Alloc",
	"BuckHashSys",
	"Frees",
	"GCCPUFraction",
	"GCSys",
	"HeapAlloc",
	"HeapIdle",
	"HeapInuse",
	"HeapObjects",
	"HeapReleased",
	"HeapSys",
	"LastGC",
	"Lookups",
	"MCacheInuse",
	"MCacheSys",
	"MSpanInuse",
	"MSpanSys",
	"Mallocs",
	"NextGC",
	"NumForcedGC",
	"NumGC",
	"OtherSys",
	"PauseTotalNs",
	"StackInuse",
	"StackSys",
	"Sys",
	"TotalAlloc",
}

// PrepareRuntimeMetrics готовит список Runtime-метрик.
func PrepareRuntimeMetrics(memStats *runtime.MemStats) []Metrics {
	data := make([]Metrics, 0)

	for _, name := range runtimeMetricsNames {
		val := reflect.ValueOf(*memStats).FieldByName(name)

		value := getFloat64FromValueObj(val)

		data = append(data, Metrics{
			ID:    name,
			MType: "gauge",
			Value: &value,
		})
	}

	randValue := rand.Float64()
	data = append(data, Metrics{
		ID:    "RandomValue",
		MType: "gauge",
		Value: &randValue,
	})

	return data
}

func getFloat64FromValueObj(val reflect.Value) float64 {
	var value float64

	if val.CanFloat() {
		value = val.Float()
	}

	if val.CanUint() {
		value = float64(val.Uint())
	}

	return value
}

// PrepareMemstatsMetrics готовит список Memstats-метрик.
func PrepareMemstatsMetrics(stats *mem.SwapMemoryStat) []Metrics {
	data := make([]Metrics, 0)

	total := float64(stats.Total)
	data = append(data, Metrics{
		ID:    "TotalMemory",
		MType: "gauge",
		Value: &total,
	})

	free := float64(stats.Free)
	data = append(data, Metrics{
		ID:    "FreeMemory",
		MType: "gauge",
		Value: &free,
	})

	return data
}

// PrepareCPUMetrics готовит список CPU-метрик.
func PrepareCPUMetrics(stats []float64) []Metrics {
	data := make([]Metrics, 0)

	for index, value := range stats {
		value := value
		data = append(data, Metrics{
			ID:    fmt.Sprintf("CPUutilization%d", index+1),
			MType: "gauge",
			Value: &value,
		})
	}

	return data
}

// PreparePollCounterMetric отдает PollCount-метрику.
func PreparePollCounterMetric(PollCount int64) Metrics {
	return Metrics{
		ID:    "PollCount",
		MType: "counter",
		Delta: &PollCount,
	}
}

// Metrics структура метрики
type Metrics struct {
	ID    string   `json:"id"`              // Название
	MType string   `json:"type"`            // Тип, возможные значения: counter, gauge
	Delta *int64   `json:"delta,omitempty"` // Значение (для типа counter)
	Value *float64 `json:"value,omitempty"` // Значение (для типа gauge)
}

func GetMetricsFromBytes(data []byte) ([]Metrics, []Metrics, error) {
	var request []Metrics

	counters := []Metrics{}
	gauges := []Metrics{}

	dec := json.NewDecoder(bytes.NewReader(data))
	if err := dec.Decode(&request); err != nil && err != io.EOF {
		return counters, gauges, err
	}

	for _, metrica := range request {
		if metrica.MType == "" || metrica.ID == "" || (metrica.Value == nil && metrica.Delta == nil) {
			continue
		}

		switch metrica.MType {
		case "counter":
			counters = append(counters, metrica)
		case "gauge":
			gauges = append(gauges, metrica)
		}
	}

	return counters, gauges, nil
}
