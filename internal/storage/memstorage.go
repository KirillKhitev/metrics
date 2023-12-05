package storage

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
