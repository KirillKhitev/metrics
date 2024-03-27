package storage

import (
	"context"
	"github.com/KirillKhitev/metrics/internal/flags"
	"github.com/KirillKhitev/metrics/internal/metrics"
	"testing"
)

const benchDBConnectionString = "host=localhost user=postgres password=sa123456 dbname=testdb sslmode=disable"

var testDBStorage = getBenchDBStorage()

func getBenchDBStorage() *DBStorage {
	flags.Args.DBConnectionString = benchDBConnectionString

	storage := &DBStorage{}
	_ = storage.Init(context.TODO())

	return storage
}

func BenchmarkDBStorage_GetCounter(b *testing.B) {
	b.Run("GetCounter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = testDBStorage.GetCounter(context.TODO(), "m2")
		}
	})
}

func BenchmarkDBStorage_GetGauge(b *testing.B) {
	b.Run("GetGauge", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = testDBStorage.GetGauge(context.TODO(), "Alloc")
		}
	})
}

func BenchmarkDBStorage_GetCounters(b *testing.B) {
	b.Run("GetCounters", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = testDBStorage.GetCounters(context.TODO(), []string{"m1", "m2"})
		}
	})
}

func BenchmarkDBStorage_GetGauges(b *testing.B) {
	b.Run("GetGauges", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = testDBStorage.GetGauges(context.TODO(), []string{"Alloc", "SomeMetric1"})
		}
	})
}

func BenchmarkDBStorage_UpdateCounter(b *testing.B) {
	b.Run("UpdateCounter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = testDBStorage.UpdateCounter(context.TODO(), "m1", 20000)
		}
	})
}

func BenchmarkDBStorage_UpdateGauge(b *testing.B) {
	b.Run("UpdateGauge", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = testDBStorage.UpdateGauge(context.TODO(), "Alloc", 1200.50)
		}
	})
}

func BenchmarkDBStorage_UpdateCounters(b *testing.B) {
	b.Run("UpdateCounters", func(b *testing.B) {
		var val1 int64 = 2000
		var val2 int64 = 500

		for i := 0; i < b.N; i++ {
			_ = testDBStorage.UpdateCounters(context.TODO(), []metrics.Metrics{
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

func BenchmarkDBStorage_UpdateGauges(b *testing.B) {
	var val1 = 2000.25
	var val2 = 500.00

	for i := 0; i < b.N; i++ {
		_ = testDBStorage.UpdateGauges(context.TODO(), []metrics.Metrics{
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

func BenchmarkDBStorage_TrySaveToFile(b *testing.B) {
	b.Run("TrySaveToFile", func(b *testing.B) {
		flags.Args.FileStoragePath = "./metrics_bench.json"
		for i := 0; i < b.N; i++ {
			_ = testDBStorage.TrySaveToFile()
		}
	})
}

func BenchmarkDBStorage_GetCounterList(b *testing.B) {
	b.Run("GetCounterList", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = testDBStorage.GetCounterList(context.TODO())
		}
	})
}

func BenchmarkDBStorage_GetGaugeList(b *testing.B) {
	b.Run("GetGaugeList", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = testDBStorage.GetGaugeList(context.TODO())
		}
	})
}
