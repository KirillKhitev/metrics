package handlers

import (
	"context"
	"github.com/KirillKhitev/metrics/internal/storage"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListHandler_ServeHTTP(t *testing.T) {
	type want struct {
		code        int
		body        string
		contentType string
	}
	type args struct {
		method  string
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
			},
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		},
		{
			name: "positive GET test #2",
			args: args{
				method: http.MethodGet,
				counter: map[string]int64{
					"PollCount": 10,
				},
				gauge: map[string]float64{
					"Alloc":      100.00,
					"SomeMetric": -324.44,
				},
			},
			want: want{
				code:        http.StatusOK,
				contentType: `text/html`,
				body:        `<b>Counter:</b><br/><p>PollCount: 10</p><br/><br/><b>Gauge:</b><br/><p>Alloc: 100</p><p>SomeMetric: -324.44</p>`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appStorage := storage.MemStorage{}

			if err := appStorage.Init(context.Background()); err != nil {
				t.Fatal("Не удалось создать хранилище")
			}

			for name, value := range tt.args.counter {
				appStorage.UpdateCounter(context.TODO(), name, value)
			}

			for name, value := range tt.args.gauge {
				appStorage.UpdateGauge(context.TODO(), name, value)
			}

			ch := &ListHandler{
				Storage: &appStorage,
			}

			request := httptest.NewRequest(tt.args.method, "/", nil)
			w := httptest.NewRecorder()

			ch.ServeHTTP(w, request)

			res := w.Result()
			defer res.Body.Close()

			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			require.Equal(t, tt.want.code, res.StatusCode)
			require.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			require.ElementsMatch(t, []byte(tt.want.body), resBody)
		})
	}
}
