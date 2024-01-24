package storage

import (
	"context"
	"github.com/KirillKhitev/metrics/internal/metrics"
)

const AttemptCount int = 4

type Repository interface {
	UpdateCounter(ctx context.Context, name string, value int64) error
	UpdateCounters(ctx context.Context, data []metrics.Metrics) error
	GetCounter(ctx context.Context, name string) (int64, error)
	GetCounters(ctx context.Context, data []string) (map[string]int64, error)

	UpdateGauge(ctx context.Context, name string, value float64) error
	UpdateGauges(ctx context.Context, data []metrics.Metrics) error
	GetGauge(ctx context.Context, name string) (float64, error)
	GetGauges(ctx context.Context, data []string) (map[string]float64, error)

	GetCounterList(ctx context.Context) map[string]int64
	GetGaugeList(ctx context.Context) map[string]float64

	Init(ctx context.Context) error
	Ping(ctx context.Context) error

	TrySaveToFile() error
	Close() error
}
