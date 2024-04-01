package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KirillKhitev/metrics/internal/metrics"
	"github.com/KirillKhitev/metrics/internal/storage"
)

func TestUpdatesHandler_ServeHTTP(t *testing.T) {
	type want struct {
		code int
		body string
	}
	type args struct {
		method  string
		metrics []metrics.Metrics
	}

	alloc := 300.55
	memory := 4354.456
	var someMetric int64 = 55
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
				metrics: []metrics.Metrics{
					{
						ID:    "Alloc",
						MType: "gauge",
						Value: &alloc,
					},
					{
						ID:    "Memory",
						MType: "gauge",
						Value: &memory,
					},
					{
						ID:    "PollCount",
						MType: "counter",
						Delta: &pollCount,
					},
					{
						ID:    "SomeMetric34",
						MType: "counter",
						Delta: &someMetric,
					},
				},
			},
			want: want{
				code: http.StatusOK,
				body: `[
				    {
				        "id": "SomeMetric34",
				        "type": "counter",
				        "delta": 55
				    },
				    {
				        "id": "PollCount",
				        "type": "counter",
				        "delta": 130
				    },
				    {
				        "id": "Alloc",
				        "type": "gauges",
				        "value": 300.55
				    },
				    {
				        "id": "Memory",
				        "type": "gauges",
				        "value": 4354.456
				    }
				]`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appStorage := storage.MemStorage{}

			if err := appStorage.Init(context.Background()); err != nil {
				t.Fatal("Не удалось создать хранилище")
			}

			if err := appStorage.UpdateCounter(context.TODO(), "PollCount", 30); err != nil {
				t.Fatal(err)
			}

			if err := appStorage.UpdateGauge(context.TODO(), "Alloc", 125.20); err != nil {
				t.Fatal(err)
			}

			ch := &UpdatesHandler{
				Storage: &appStorage,
			}

			str, _ := json.Marshal(tt.args.metrics)

			buf := bytes.NewBuffer(str)

			request := httptest.NewRequest(tt.args.method, "/updates/", buf)
			w := httptest.NewRecorder()

			ch.ServeHTTP(w, request)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.code, res.StatusCode)

			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			require.Equal(t, tt.want.code, res.StatusCode)

			metricsResponse := []metrics.Metrics{}
			metricsWant := []metrics.Metrics{}

			if err := json.Unmarshal(resBody, &metricsResponse); err != nil {
				t.Fatal(err)
			}

			if err := json.Unmarshal([]byte(tt.want.body), &metricsWant); err != nil {
				t.Fatal(err)
			}

			require.ElementsMatch(t, metricsWant, metricsResponse)
		})
	}
}
