package metrics

import (
	"fmt"
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

func PrepareGaugeForSend(memStats *runtime.MemStats) (gauge map[string]string) {
	gauge = make(map[string]string)

	for _, name := range runtimeMetricsNames {
		val := reflect.ValueOf(*memStats).FieldByName(name)

		var value float64

		if val.CanFloat() {
			value = val.Float()
		}

		if val.CanUint() {
			value = float64(val.Uint())
		}

		gauge[name] = fmt.Sprintf("%f", value)
	}

	gauge["RandomValue"] = fmt.Sprintf("%f", rand.Float64())

	return
}

func PrepareCounterForSend(PollCount int64) (counter map[string]string) {
	counter = make(map[string]string)
	counter["PollCount"] = fmt.Sprintf("%d", PollCount)

	return
}
