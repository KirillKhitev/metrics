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
		wantCounter map[string]int64
	}{
		{
			name: "positive test #1",
			args: args{
				PollCount: 10,
			},
			wantCounter: map[string]int64{
				"PollCount": 10,
			},
		},
		{
			name: "positive test #2",
			args: args{
				PollCount: -10,
			},
			wantCounter: map[string]int64{
				"PollCount": -10,
			},
		},
		{
			name: "positive test #3",
			args: args{
				PollCount: 234234234,
			},
			wantCounter: map[string]int64{
				"PollCount": 234234234,
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

	stats1 := new(runtime.MemStats)

	*stats1 = runtime.MemStats{
		Alloc:         1000,
		BuckHashSys:   2000,
		GCCPUFraction: 10123.45,
	}

	stats2 := new(runtime.MemStats)

	*stats2 = runtime.MemStats{
		Alloc:         10000099,
		BuckHashSys:   0,
		Frees:         1000000,
		GCCPUFraction: -10123456.45,
	}

	tests := []struct {
		name      string
		args      args
		wantGauge map[string]float64
	}{
		{
			name: "positive test #1",
			args: args{
				memStats: stats1,
			},
			wantGauge: map[string]float64{
				"Alloc":         1000,
				"BuckHashSys":   2000,
				"Frees":         0,
				"GCCPUFraction": 10123.45,
				"GCSys":         0,
				"HeapAlloc":     0,
				"HeapIdle":      0,
				"HeapInuse":     0,
				"HeapObjects":   0,
				"HeapReleased":  0,
				"HeapSys":       0,
				"LastGC":        0,
				"Lookups":       0,
				"MCacheInuse":   0,
				"MCacheSys":     0,
				"MSpanInuse":    0,
				"MSpanSys":      0,
				"Mallocs":       0,
				"NextGC":        0,
				"NumForcedGC":   0,
				"NumGC":         0,
				"OtherSys":      0,
				"PauseTotalNs":  0,
				"StackInuse":    0,
				"StackSys":      0,
				"Sys":           0,
				"TotalAlloc":    0,
				"RandomValue":   0,
			},
		},
		{
			name: "positive test #2",
			args: args{
				memStats: stats2,
			},
			wantGauge: map[string]float64{
				"Alloc":         10000099,
				"BuckHashSys":   0,
				"Frees":         1000000,
				"GCCPUFraction": -10123456.45,
				"GCSys":         0,
				"HeapAlloc":     0,
				"HeapIdle":      0,
				"HeapInuse":     0,
				"HeapObjects":   0,
				"HeapReleased":  0,
				"HeapSys":       0,
				"LastGC":        0,
				"Lookups":       0,
				"MCacheInuse":   0,
				"MCacheSys":     0,
				"MSpanInuse":    0,
				"MSpanSys":      0,
				"Mallocs":       0,
				"NextGC":        0,
				"NumForcedGC":   0,
				"NumGC":         0,
				"OtherSys":      0,
				"PauseTotalNs":  0,
				"StackInuse":    0,
				"StackSys":      0,
				"Sys":           0,
				"TotalAlloc":    0,
				"RandomValue":   0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotGauge := PrepareGaugeForSend(tt.args.memStats)
			gotGauge["RandomValue"] = 0

			if !reflect.DeepEqual(gotGauge, tt.wantGauge) {
				t.Errorf("PrepareGaugeForSend() = %v, want %v", gotGauge, tt.wantGauge)
			}
		})
	}
}
