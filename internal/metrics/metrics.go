package metrics

import (
	"math/rand"
	"reflect"
	"runtime"
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

func PrepareGaugeForSend(memStats *runtime.MemStats) (gauge map[string]float64) {
	gauge = make(map[string]float64)

	for _, name := range runtimeMetricsNames {
		val := reflect.ValueOf(*memStats).FieldByName(name)

		var value float64

		if val.CanFloat() {
			value = val.Float()
		}

		if val.CanUint() {
			value = float64(val.Uint())
		}

		gauge[name] = value
	}

	gauge["RandomValue"] = rand.Float64()

	return
}

func PrepareCounterForSend(PollCount int64) (counter map[string]int64) {
	counter = make(map[string]int64)
	counter["PollCount"] = PollCount

	return
}

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}
