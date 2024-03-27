package storage

import (
	"context"
	"github.com/KirillKhitev/metrics/internal/flags"
	"github.com/KirillKhitev/metrics/internal/metrics"
	"testing"
)

var testMemStorage = &MemStorage{
	Counter: map[string]int64{
		"m1": 100,
		"m2": 100000,
		"m3": -100000,
	},
	Gauge: map[string]float64{
		"Alloc":       10000.00,
		"SomeMetric1": 10.00,
		"SomeMetric2": -10000,
	},
}

func BenchmarkMemStorage_GetCounter(b *testing.B) {
	b.Run("GetCounter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = testMemStorage.GetCounter(context.TODO(), "m2")
		}
	})
}

func BenchmarkMemStorage_GetGauge(b *testing.B) {
	b.Run("GetGauge", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = testMemStorage.GetGauge(context.TODO(), "Alloc")
		}
	})
}

func BenchmarkMemStorage_GetCounters(b *testing.B) {
	b.Run("GetCounters", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = testMemStorage.GetCounters(context.TODO(), []string{"m1", "m2"})
		}
	})
}

func BenchmarkMemStorage_GetGauges(b *testing.B) {
	b.Run("GetGauges", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = testMemStorage.GetGauges(context.TODO(), []string{"Alloc", "SomeMetric1"})
		}
	})
}

func BenchmarkMemStorage_UpdateCounter(b *testing.B) {
	b.Run("UpdateCounter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = testMemStorage.UpdateCounter(context.TODO(), "m1", 20000)
		}
	})
}

func BenchmarkMemStorage_UpdateGauge(b *testing.B) {
	b.Run("UpdateGauge", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = testMemStorage.UpdateGauge(context.TODO(), "Alloc", 1200.50)
		}
	})
}

func BenchmarkMemStorage_UpdateCounters(b *testing.B) {
	b.Run("UpdateCounters", func(b *testing.B) {
		var val1 int64 = 2000
		var val2 int64 = 500

		for i := 0; i < b.N; i++ {
			_ = testMemStorage.UpdateCounters(context.TODO(), []metrics.Metrics{
				{
					ID:    "m1",
					MType: "counter",
					Delta: &val1,
				},
				{
					ID:    "m2",
					MType: "counter",
					Delta: &val2,
				},
			})
		}
	})
}

func BenchmarkMemStorage_UpdateGauges(b *testing.B) {
	var val1 = 2000.25
	var val2 = 500.00

	for i := 0; i < b.N; i++ {
		_ = testMemStorage.UpdateGauges(context.TODO(), []metrics.Metrics{
			{
				ID:    "Alloc",
				MType: "gauge",
				Value: &val1,
			},
			{
				ID:    "SomeMetric1",
				MType: "gauge",
				Value: &val2,
			},
		})
	}
}

func BenchmarkMemStorage_TrySaveToFile(b *testing.B) {
	b.Run("TrySaveToFile", func(b *testing.B) {
		flags.Args.FileStoragePath = "./metrics_bench.json"
		for i := 0; i < b.N; i++ {
			_ = testMemStorage.TrySaveToFile()
		}
	})
}

func BenchmarkMemStorage_GetCounterList(b *testing.B) {
	b.Run("GetCounterList", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = testMemStorage.GetCounterList(context.TODO())
		}
	})
}

func BenchmarkMemStorage_GetGaugeList(b *testing.B) {
	b.Run("GetGaugeList", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = testMemStorage.GetGaugeList(context.TODO())
		}
	})
}
