package client

import (
	"context"

	"github.com/KirillKhitev/metrics/internal/metrics"
)

type Client interface {
	Send(ctx context.Context, data []metrics.Metrics) error
	Close() error
}
