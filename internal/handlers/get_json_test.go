package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/KirillKhitev/metrics/internal/metrics"
	"github.com/KirillKhitev/metrics/internal/storage"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetJSONHandler_ServeHTTP(t *testing.T) {
	type want struct {
		code int
		body string
	}
	type args struct {
		method string
		metric metrics.Metrics
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "positive Post test #1",
			args: args{
				method: http.MethodPost,
				metric: metrics.Metrics{
					ID:    "Alloc",
					MType: "gauge",
				},
			},
			want: want{
				code: http.StatusOK,
				body: `{"id":"Alloc", "type":"gauge", "value":3000.555}`,
			},
		},
		{
			name: "positive Post test #2",
			args: args{
				method: http.MethodPost,
				metric: metrics.Metrics{
					ID:    "PollCount",
					MType: "counter",
				},
			},
			want: want{
				code: http.StatusOK,
				body: `{"id":"PollCount", "type":"counter", "delta":100}`,
			},
		},
		{
			name: "negative Post test #3",
			args: args{
				method: http.MethodPost,
				metric: metrics.Metrics{
					ID:    "SomeMetric",
					MType: "counter",
				},
			},
			want: want{
				code: http.StatusNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appStorage := storage.MemStorage{}

			if err := appStorage.Init(context.Background()); err != nil {
				t.Fatal("Не удалось создать хранилище")
			}

			_ = appStorage.UpdateCounter(context.TODO(), "PollCount", 100)
			_ = appStorage.UpdateGauge(context.TODO(), "Alloc", 3000.555)

			ch := &GetJSONHandler{
				Storage: &appStorage,
			}

			str, _ := json.Marshal(tt.args.metric)

			buf := bytes.NewBuffer(str)

			request := httptest.NewRequest(tt.args.method, "/value/", buf)
			w := httptest.NewRecorder()

			ch.ServeHTTP(w, request)

			res := w.Result()
			defer res.Body.Close()

			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			require.Equal(t, tt.want.code, res.StatusCode)

			if tt.want.body != "" {
				require.JSONEq(t, tt.want.body, string(resBody))
			}
		})
	}
}
