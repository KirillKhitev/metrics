package storage

import "errors"

type MemStorage struct {
	counter map[string]int64
	gauge   map[string]float64
}

func (s *MemStorage) UpdateCounter(name string, value int64) error {
	s.counter[name] += value

	return nil
}

func (s *MemStorage) UpdateGauge(name string, value float64) error {
	s.gauge[name] = value

	return nil
}

func (s *MemStorage) Init() error {
	s.counter = make(map[string]int64)
	s.gauge = make(map[string]float64)

	return nil
}

var ErrNotFound = errors.New("not found")

func (s *MemStorage) GetCounter(name string) (v int64, err error) {
	v, ok := s.counter[name]

	if !ok {
		return v, ErrNotFound
	}

	return v, nil
}

func (s *MemStorage) GetGauge(name string) (v float64, err error) {
	v, ok := s.gauge[name]

	if !ok {
		return v, ErrNotFound
	}

	return v, nil
}

func (s *MemStorage) GetCounterList() map[string]int64 {
	return s.counter
}

func (s *MemStorage) GetGaugeList() map[string]float64 {
	return s.gauge
}
