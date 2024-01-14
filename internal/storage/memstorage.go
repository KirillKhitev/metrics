package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/KirillKhitev/metrics/internal/flags"
	"github.com/KirillKhitev/metrics/internal/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"io"
	"os"
)

type MemStorage struct {
	DB      *sql.DB
	Counter map[string]int64
	Gauge   map[string]float64
}

func (s *MemStorage) UpdateCounter(name string, value int64) error {
	s.Counter[name] += value

	return nil
}

func (s *MemStorage) UpdateGauge(name string, value float64) error {
	s.Gauge[name] = value

	return nil
}

func (s *MemStorage) Init() error {
	s.Counter = make(map[string]int64)
	s.Gauge = make(map[string]float64)

	if err := s.initDBConnect(); err != nil {
		return err
	}

	if !flags.Args.Restore {
		return nil
	}

	file, err := os.OpenFile(flags.Args.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Log.Error("Error by get metrics from file", zap.Error(err))
		return nil
	}

	defer file.Close()

	if err := json.NewDecoder(file).Decode(s); err != nil && err != io.EOF {
		logger.Log.Error("Error by decode metrics from json", zap.Error(err))
		return nil
	}

	return nil
}

var ErrNotFound = errors.New("not found")

func (s *MemStorage) GetCounter(name string) (v int64, err error) {
	v, ok := s.Counter[name]

	if !ok {
		return v, ErrNotFound
	}

	return v, nil
}

func (s *MemStorage) GetGauge(name string) (v float64, err error) {
	v, ok := s.Gauge[name]

	if !ok {
		return v, ErrNotFound
	}

	return v, nil
}

func (s *MemStorage) GetCounterList() map[string]int64 {
	return s.Counter
}

func (s *MemStorage) GetGaugeList() map[string]float64 {
	return s.Gauge
}

func (s *MemStorage) SaveToFile() error {
	logger.Log.Info("Сохраняем метрики в файл")

	file, err := os.OpenFile(flags.Args.FileStoragePath, os.O_WRONLY|os.O_CREATE, 0666)
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

func (s *MemStorage) initDBConnect() error {
	db, err := sql.Open("pgx", flags.Args.DbConnectionString)
	if err != nil {
		return err
	}

	s.DB = db

	return nil
}

func (s *MemStorage) CloseDB() {
	logger.Log.Info("Close connect to DB")
	s.DB.Close()
}
