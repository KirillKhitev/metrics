package storage

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMemStorage_GetCounter(t *testing.T) {
	type storage struct {
		counter map[string]int64
		gauge   map[string]float64
	}
	type args struct {
		name string
	}

	storageApp := storage{
		counter: map[string]int64{
			"m1": 100,
			"m2": 100000,
			"m3": -100000,
		},
		gauge: map[string]float64{
			"Alloc":       10000.00,
			"SomeMetric1": 10.00,
			"SomeMetric2": -10000,
		},
	}

	tests := []struct {
		name    string
		storage storage
		args    args
		want    int64
		wantErr bool
	}{
		{
			name:    "positive get counter test #1",
			storage: storageApp,
			args: args{
				name: "m2",
			},
			want:    100000,
			wantErr: false,
		},
		{
			name:    "positive get counter test #2",
			storage: storageApp,
			args: args{
				name: "m3",
			},
			want:    -100000,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				Counter: tt.storage.counter,
				Gauge:   tt.storage.gauge,
			}
			got, err := s.GetCounter(context.TODO(), tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetCounter() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_GetGauge(t *testing.T) {
	type storage struct {
		counter map[string]int64
		gauge   map[string]float64
	}
	type args struct {
		name string
	}

	storageApp := storage{
		counter: map[string]int64{
			"m1": 100,
			"m2": 100000,
			"m3": -100000,
		},
		gauge: map[string]float64{
			"Alloc":       10000.00,
			"SomeMetric1": 10.00,
			"SomeMetric2": -10000,
		},
	}

	tests := []struct {
		name    string
		storage storage
		args    args
		want    float64
		wantErr bool
	}{
		{
			name:    "positive get gauge test #1",
			storage: storageApp,
			args: args{
				name: "Alloc",
			},
			want:    10000.00,
			wantErr: false,
		},
		{
			name:    "positive get gauge test #2",
			storage: storageApp,
			args: args{
				name: "SomeMetric1",
			},
			want:    10.00,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				Counter: tt.storage.counter,
				Gauge:   tt.storage.gauge,
			}
			got, err := s.GetGauge(context.TODO(), tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGauge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetGauge() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_Init(t *testing.T) {
	type storage struct {
		counter map[string]int64
		gauge   map[string]float64
	}

	storageApp := storage{
		counter: map[string]int64{},
		gauge:   map[string]float64{},
	}

	tests := []struct {
		name    string
		storage storage
		wantErr bool
	}{
		{
			name:    "positive init storage test #1",
			storage: storageApp,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{}
			if err := s.Init(); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemStorage_UpdateCounter(t *testing.T) {
	type storage struct {
		counter map[string]int64
		gauge   map[string]float64
	}

	storageApp := storage{
		counter: map[string]int64{
			"m1": 100,
			"m2": 100000,
			"m3": -100000,
		},
		gauge: map[string]float64{
			"Alloc":       10000.00,
			"SomeMetric1": 10.00,
			"SomeMetric2": -10000,
		},
	}

	type args struct {
		name  string
		value int64
	}

	tests := []struct {
		name      string
		storage   storage
		args      args
		wantValue int64
		wantErr   bool
	}{
		{
			name:    "positive update counter test #1",
			storage: storageApp,
			args: args{
				name:  "m2",
				value: 20,
			},
			wantValue: 100020.00,
			wantErr:   false,
		},
		{
			name:    "positive update counter test #2",
			storage: storageApp,
			args: args{
				name:  "m4",
				value: 200,
			},
			wantValue: 200,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				Counter: tt.storage.counter,
				Gauge:   tt.storage.gauge,
			}

			err := s.UpdateCounter(context.TODO(), tt.args.name, tt.args.value)

			if tt.wantErr {
				require.NotNil(t, err)
			}

			val, _ := s.GetCounter(context.TODO(), tt.args.name)

			require.Equal(t, tt.wantValue, val)
		})
	}
}

func TestMemStorage_UpdateGauge(t *testing.T) {
	type storage struct {
		counter map[string]int64
		gauge   map[string]float64
	}

	storageApp := storage{
		counter: map[string]int64{
			"m1": 100,
			"m2": 100000,
			"m3": -100000,
		},
		gauge: map[string]float64{
			"Alloc":       10000.00,
			"SomeMetric1": 10.00,
			"SomeMetric2": -10000,
		},
	}

	type args struct {
		name  string
		value float64
	}

	tests := []struct {
		name      string
		storage   storage
		args      args
		wantValue float64
		wantErr   bool
	}{
		{
			name:    "positive update gauge test #1",
			storage: storageApp,
			args: args{
				name:  "Alloc",
				value: 1000,
			},
			wantValue: 1000.00,
			wantErr:   false,
		},
		{
			name:    "positive update gauge test #2",
			storage: storageApp,
			args: args{
				name:  "SomeMetric2",
				value: -5000.55,
			},
			wantValue: -5000.55,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				Counter: tt.storage.counter,
				Gauge:   tt.storage.gauge,
			}

			err := s.UpdateGauge(context.TODO(), tt.args.name, tt.args.value)

			if tt.wantErr {
				require.NotNil(t, err)
			}

			val, _ := s.GetGauge(context.TODO(), tt.args.name)

			require.Equal(t, tt.wantValue, val)
		})
	}
}
