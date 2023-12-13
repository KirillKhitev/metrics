package metrics

import (
	"reflect"
	"runtime"
	"testing"
)

func TestPrepareCounterForSend(t *testing.T) {
	type args struct {
		PollCount int64
	}
	tests := []struct {
		name        string
		args        args
		wantCounter map[string]string
	}{
		{
			name: "positive test #1",
			args: args{
				PollCount: 10,
			},
			wantCounter: map[string]string{
				"PollCount": "10",
			},
		},
		{
			name: "positive test #2",
			args: args{
				PollCount: -10,
			},
			wantCounter: map[string]string{
				"PollCount": "-10",
			},
		},
		{
			name: "positive test #3",
			args: args{
				PollCount: 234234234,
			},
			wantCounter: map[string]string{
				"PollCount": "234234234",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCounter := PrepareCounterForSend(tt.args.PollCount); !reflect.DeepEqual(gotCounter, tt.wantCounter) {
				t.Errorf("PrepareCounterForSend() = %v, want %v", gotCounter, tt.wantCounter)
			}
		})
	}
}

func TestPrepareGaugeForSend(t *testing.T) {
	type args struct {
		memStats *runtime.MemStats
	}

	tests := []struct {
		name      string
		args      args
		wantGauge map[string]string
	}{
		{
			name: "positive test #1",
			args: args{
				memStats: &runtime.MemStats{
					Alloc:         1000,
					BuckHashSys:   2000,
					GCCPUFraction: 10123.45,
				},
			},
			wantGauge: map[string]string{
				"Alloc":         "1000.000000",
				"BuckHashSys":   "2000.000000",
				"Frees":         "0.000000",
				"GCCPUFraction": "10123.450000",
				"GCSys":         "0.000000",
				"HeapAlloc":     "0.000000",
				"HeapIdle":      "0.000000",
				"HeapInuse":     "0.000000",
				"HeapObjects":   "0.000000",
				"HeapReleased":  "0.000000",
				"HeapSys":       "0.000000",
				"LastGC":        "0.000000",
				"Lookups":       "0.000000",
				"MCacheInuse":   "0.000000",
				"MCacheSys":     "0.000000",
				"MSpanInuse":    "0.000000",
				"MSpanSys":      "0.000000",
				"Mallocs":       "0.000000",
				"NextGC":        "0.000000",
				"NumForcedGC":   "0.000000",
				"NumGC":         "0.000000",
				"OtherSys":      "0.000000",
				"PauseTotalNs":  "0.000000",
				"StackInuse":    "0.000000",
				"StackSys":      "0.000000",
				"Sys":           "0.000000",
				"TotalAlloc":    "0.000000",
				"RandomValue":   "0.000000",
			},
		},
		{
			name: "positive test #2",
			args: args{
				memStats: &runtime.MemStats{
					Alloc:         10000099,
					BuckHashSys:   0,
					Frees:         1000000,
					GCCPUFraction: -10123456.45,
				},
			},
			wantGauge: map[string]string{
				"Alloc":         "10000099.000000",
				"BuckHashSys":   "0.000000",
				"Frees":         "1000000.000000",
				"GCCPUFraction": "-10123456.450000",
				"GCSys":         "0.000000",
				"HeapAlloc":     "0.000000",
				"HeapIdle":      "0.000000",
				"HeapInuse":     "0.000000",
				"HeapObjects":   "0.000000",
				"HeapReleased":  "0.000000",
				"HeapSys":       "0.000000",
				"LastGC":        "0.000000",
				"Lookups":       "0.000000",
				"MCacheInuse":   "0.000000",
				"MCacheSys":     "0.000000",
				"MSpanInuse":    "0.000000",
				"MSpanSys":      "0.000000",
				"Mallocs":       "0.000000",
				"NextGC":        "0.000000",
				"NumForcedGC":   "0.000000",
				"NumGC":         "0.000000",
				"OtherSys":      "0.000000",
				"PauseTotalNs":  "0.000000",
				"StackInuse":    "0.000000",
				"StackSys":      "0.000000",
				"Sys":           "0.000000",
				"TotalAlloc":    "0.000000",
				"RandomValue":   "0.000000",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotGauge := PrepareGaugeForSend(tt.args.memStats)
			gotGauge["RandomValue"] = "0.000000"

			if !reflect.DeepEqual(gotGauge, tt.wantGauge) {
				t.Errorf("PrepareGaugeForSend() = %v,\r\n want %v", gotGauge, tt.wantGauge)
			}
		})
	}
}
