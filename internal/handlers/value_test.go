package handlers

import (
	"github.com/KirillKhitev/metrics/internal/storage"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValueHandler_ServeHTTP(t *testing.T) {
	type want struct {
		code int
		body string
	}
	type args struct {
		method  string
		addr    string
		counter map[string]int64
		gauge   map[string]float64
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "negative POST test #1",
			args: args{
				method: http.MethodPost,
				addr:   "/value/",
			},
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		},
		{
			name: "negative GET test #2",
			args: args{
				method: http.MethodGet,
				addr:   "/value/",
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "negative GET test #3",
			args: args{
				method: http.MethodGet,
				addr:   "/value/wrong_type/",
			},
			want: want{
				code: http.StatusNotFound,
			},
		},
		{
			name: "negative GET test #4",
			args: args{
				method: http.MethodGet,
				addr:   "/value/wrong_type/Alloc",
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "negative GET test #4",
			args: args{
				method: http.MethodGet,
				addr:   "/value/counter/",
			},
			want: want{
				code: http.StatusNotFound,
			},
		},
		{
			name: "negative GET test #5",
			args: args{
				method: http.MethodGet,
				addr:   "/value/counter/WrongMetric",
			},
			want: want{
				code: http.StatusNotFound,
			},
		},
		{
			name: "positive GET test #6",
			args: args{
				method: http.MethodGet,
				addr:   "/value/counter/PollCount",
				counter: map[string]int64{
					"PollCount": 10,
				},
				gauge: map[string]float64{
					"Alloc":      100.00,
					"SomeMetric": -324.44,
				},
			},
			want: want{
				code: http.StatusOK,
				body: "10",
			},
		},
		{
			name: "positive GET test #7",
			args: args{
				method: http.MethodGet,
				addr:   "/value/gauge/Alloc",
				counter: map[string]int64{
					"PollCount": 10,
				},
				gauge: map[string]float64{
					"Alloc":      100.00,
					"SomeMetric": -324.44,
				},
			},
			want: want{
				code: http.StatusOK,
				body: "100",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appStorage := storage.MemStorage{}

			if err := appStorage.Init(); err != nil {
				t.Fatal("Не удалось создать хранилище")
			}

			for name, value := range tt.args.counter {
				appStorage.UpdateCounter(name, value)
			}

			for name, value := range tt.args.gauge {
				appStorage.UpdateGauge(name, value)
			}

			ch := &ValueHandler{
				Storage: appStorage,
			}

			request := httptest.NewRequest(tt.args.method, tt.args.addr, nil)
			w := httptest.NewRecorder()

			ch.ServeHTTP(w, request)

			res := w.Result()
			defer res.Body.Close()

			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			require.Equal(t, tt.want.code, res.StatusCode)
			require.Equal(t, tt.want.body, string(resBody))
		})
	}
}
