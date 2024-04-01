package storage

import (
	"context"

	"github.com/KirillKhitev/metrics/internal/metrics"
)

// AttemptCount - количество попыток выполнения запроса к Хранилищу
const AttemptCount int = 4

// Repository - интерфейс, описывающий методы обновления данных в хранилище
type Repository interface {
	// UpdateCounter обновляет метрику типа counter новым значением.
	// К старому значению добавляется новое.
	UpdateCounter(ctx context.Context, name string, value int64) error

	// UpdateCounters обновляет список метрик типа counter новыми значениями.
	// К старому значению добавляется новое.
	UpdateCounters(ctx context.Context, data []metrics.Metrics) error

	// GetCounter получает значение метрики типа counter.
	GetCounter(ctx context.Context, name string) (int64, error)

	// GetCounters получает список метрик типа counter.
	GetCounters(ctx context.Context, data []string) (map[string]int64, error)

	// UpdateGauge обновляет метрику типа gauge новым значением.
	// Старое значение заменяется новым.
	UpdateGauge(ctx context.Context, name string, value float64) error

	// UpdateGauges обновляет список метрик типа gauge новыми значениями.
	// Старое значение заменяется новым.
	UpdateGauges(ctx context.Context, data []metrics.Metrics) error

	// GetGauge получает значение метрики типа gauge.
	GetGauge(ctx context.Context, name string) (float64, error)

	// GetGauges получает список метрик типа gauge.
	GetGauges(ctx context.Context, data []string) (map[string]float64, error)

	// GetCounterList получает список всех метрик типа counter.
	GetCounterList(ctx context.Context) map[string]int64

	// GetGaugeList получает список всех метрик типа gauge.
	GetGaugeList(ctx context.Context) map[string]float64

	// Инициализация хранилища
	Init(ctx context.Context) error

	// Проверка соединения с хранилищем
	Ping(ctx context.Context) error

	// Сохранение метрик в файл
	TrySaveToFile() error

	// Закрытие хранилища
	Close() error
}
