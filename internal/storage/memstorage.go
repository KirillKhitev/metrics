package storage

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/KirillKhitev/metrics/internal/flags"
	"github.com/KirillKhitev/metrics/internal/logger"
	"github.com/KirillKhitev/metrics/internal/metrics"
	_ "github.com/jackc/pgx/v5/stdlib"
	"io"
	"os"
	"time"
)

type MemStorage struct {
	Counter map[string]int64
	Gauge   map[string]float64
}

func (s *MemStorage) UpdateCounter(ctx context.Context, name string, value int64) error {
	s.Counter[name] += value

	return nil
}

func (s *MemStorage) UpdateCounters(ctx context.Context, data []metrics.Metrics) error {
	for _, metrica := range data {
		s.Counter[metrica.ID] += *metrica.Delta
	}

	return nil
}

func (s *MemStorage) UpdateGauge(ctx context.Context, name string, value float64) error {
	s.Gauge[name] = value

	return nil
}

func (s *MemStorage) UpdateGauges(ctx context.Context, data []metrics.Metrics) error {
	for _, metrica := range data {
		s.Gauge[metrica.ID] = *metrica.Value
	}

	return nil
}

func (s *MemStorage) Init() error {
	s.Counter = make(map[string]int64)
	s.Gauge = make(map[string]float64)

	if !flags.Args.Restore {
		return nil
	}

	if err := s.getDataFromFile(); err != nil {
		return err
	}

	return nil
}

func (s *MemStorage) getDataFromFile() error {
	var file *os.File
	var err error

	for i := 1; i <= ATTEMPT_COUNT; i++ {
		file, err = os.OpenFile(flags.Args.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
		if err != nil && i != ATTEMPT_COUNT {
			continue
		}

		break
	}

	if err != nil {
		return err
	}

	defer file.Close()

	if err := json.NewDecoder(file).Decode(s); err != nil && err != io.EOF {
		return err
	}

	return nil
}

var ErrNotFound = errors.New("not found")

func (s *MemStorage) GetCounter(ctx context.Context, name string) (v int64, err error) {
	v, ok := s.Counter[name]

	if !ok {
		return v, ErrNotFound
	}

	return v, nil
}

func (s *MemStorage) GetCounters(ctx context.Context, data []string) (map[string]int64, error) {
	result := map[string]int64{}

	for _, name := range data {
		v, ok := s.Counter[name]

		if !ok {
			continue
		}

		result[name] = v
	}

	return result, nil
}

func (s *MemStorage) GetGauge(ctx context.Context, name string) (v float64, err error) {
	v, ok := s.Gauge[name]

	if !ok {
		return v, ErrNotFound
	}

	return v, nil
}

func (s *MemStorage) GetGauges(ctx context.Context, data []string) (map[string]float64, error) {
	result := map[string]float64{}

	for _, name := range data {
		v, ok := s.Gauge[name]

		if !ok {
			continue
		}

		result[name] = v
	}

	return result, nil
}

func (s *MemStorage) GetCounterList(ctx context.Context) map[string]int64 {
	return s.Counter
}

func (s *MemStorage) GetGaugeList(ctx context.Context) map[string]float64 {
	return s.Gauge
}

func (s *MemStorage) Close() error {
	return nil
}

func (s *MemStorage) TrySaveToFile() error {
	logger.Log.Info("Сохраняем метрики в файл")

	var file *os.File
	var err error

	for i := 1; i <= ATTEMPT_COUNT; i++ {
		file, err = os.OpenFile(flags.Args.FileStoragePath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil && i != ATTEMPT_COUNT {
			time.Sleep(time.Duration(2*i-1) * time.Second)
			continue
		}

		break
	}

	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	if err := encoder.Encode(s); err != nil {
		return err
	}

	return nil
}

func (s *MemStorage) Ping(ctx context.Context) error {
	return errors.New("it is not DB storage")
}
