package metrics

import (
	"fmt"
	"math/rand"
	"reflect"
	"runtime"

	"github.com/shirou/gopsutil/v3/mem"
)

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

func PreparePollCounterMetric(PollCount int64) Metrics {
	return Metrics{
		ID:    "PollCount",
		MType: "counter",
		Delta: &PollCount,
	}
}

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}
