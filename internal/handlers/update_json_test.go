package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/KirillKhitev/metrics/internal/metrics"
	"github.com/KirillKhitev/metrics/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateJsonHandler_ServeHTTP(t *testing.T) {
	type want struct {
		code int
		body string
	}
	type args struct {
		method string
		metric metrics.Metrics
	}

	var alloc float64 = 300.55
	var someMetric float64 = 55.7
	var pollCount int64 = 100

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "positive test gauge metric #1",
			args: args{
				method: http.MethodPost,
				metric: metrics.Metrics{
					ID:    "Alloc",
					MType: "gauge",
					Value: &alloc,
				},
			},
			want: want{
				code: http.StatusOK,
				body: `{"id":"Alloc", "type":"gauge", "value":300.55}`,
			},
		},
		{
			name: "positive test counter metric #2",
			args: args{
				method: http.MethodPost,
				metric: metrics.Metrics{
					ID:    "PollCount",
					MType: "counter",
					Delta: &pollCount,
				},
			},
			want: want{
				code: http.StatusOK,
				body: `{"id":"PollCount", "type":"counter", "delta":100}`,
			},
		},
		{
			name: "positive test gauge metric #3",
			args: args{
				method: http.MethodPost,
				metric: metrics.Metrics{
					ID:    "SomeMetric",
					MType: "gauge",
					Value: &someMetric,
				},
			},
			want: want{
				code: http.StatusOK,
				body: `{"id":"SomeMetric", "type":"gauge", "value":55.7}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appStorage := storage.MemStorage{}

			if err := appStorage.Init(); err != nil {
				t.Fatal("Не удалось создать хранилище")
			}

			appStorage.UpdateCounter("PollCounter", 30)
			appStorage.UpdateGauge("Alloc", 125.20)

			ch := &UpdateJsonHandler{
				Storage: appStorage,
			}

			str, _ := json.Marshal(tt.args.metric)

			buf := bytes.NewBuffer(str)

			request := httptest.NewRequest(tt.args.method, "/update/", buf)
			w := httptest.NewRecorder()

			ch.ServeHTTP(w, request)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.code, res.StatusCode)

			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			require.Equal(t, tt.want.code, res.StatusCode)
			require.JSONEq(t, tt.want.body, string(resBody))
		})
	}
}
