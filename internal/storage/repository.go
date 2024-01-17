package storage

import (
	"context"
)

type Repository interface {
	UpdateCounter(ctx context.Context, name string, value int64) error
	GetCounter(ctx context.Context, name string) (int64, error)

	UpdateGauge(ctx context.Context, name string, value float64) error
	GetGauge(ctx context.Context, name string) (float64, error)

	GetCounterList(ctx context.Context) map[string]int64
	GetGaugeList(ctx context.Context) map[string]float64

	Init() error
	Ping(ctx context.Context) error

	TrySaveToFile() error
	Close() error
}
