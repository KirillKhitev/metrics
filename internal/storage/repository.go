package storage

type Repository interface {
	UpdateCounter(name string, value int64) error
	UpdateGauge(name string, value float64) error
	Init() error
}
