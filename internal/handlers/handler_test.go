package handlers

import (
	"github.com/KirillKhitev/metrics/internal/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateHandler_ServeHTTP(t *testing.T) {
	type want struct {
		code         int
		counterValue int64
		gaugeValue   float64
	}
	type args struct {
		method       string
		addr         string
		counterValue int64
		gaugeValue   float64
		name         string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "negative test GET #1",
			args: args{
				method: http.MethodGet,
				addr:   `/update/`,
			},
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		},
		{
			name: "negative test empty type metric #2",
			args: args{
				method: http.MethodGet,
				addr:   `/update/`,
			},
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		},
		{
			name: "negative test wrong type metric #3",
			args: args{
				method: http.MethodPost,
				addr:   `/update/wrong_type/Alloc/10`,
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "negative test empty name metric #4",
			args: args{
				method: http.MethodPost,
				addr:   `/update/counter/`,
			},
			want: want{
				code: http.StatusNotFound,
			},
		},
		{
			name: "negative test empty value metric #5",
			args: args{
				method: http.MethodPost,
				addr:   `/update/counter/Alloc/`,
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "positive test counter metric #6",
			args: args{
				method:       http.MethodPost,
				addr:         `/update/counter/Alloc/10`,
				counterValue: 15,
				gaugeValue:   10.00,
				name:         "Alloc",
			},
			want: want{
				code:         http.StatusOK,
				counterValue: 25,
				gaugeValue:   10.00,
			},
		},
		{
			name: "positive test counter metric #7",
			args: args{
				method:       http.MethodPost,
				addr:         `/update/counter/Alloc/-10`,
				counterValue: 15,
				gaugeValue:   10.00,
				name:         "Alloc",
			},
			want: want{
				code:         http.StatusOK,
				counterValue: 5,
				gaugeValue:   10.00,
			},
		},
		{
			name: "positive test gauge metric #8",
			args: args{
				method:       http.MethodPost,
				addr:         `/update/gauge/Alloc/15.50`,
				counterValue: 15,
				gaugeValue:   10.00,
				name:         "Alloc",
			},
			want: want{
				code:         http.StatusOK,
				counterValue: 15,
				gaugeValue:   15.50,
			},
		},
		{
			name: "positive test gauge metric #9",
			args: args{
				method:       http.MethodPost,
				addr:         `/update/gauge/Alloc/-15.50`,
				counterValue: 15,
				gaugeValue:   10.00,
				name:         "Alloc",
			},
			want: want{
				code:         http.StatusOK,
				counterValue: 15,
				gaugeValue:   -15.50,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appStorage := storage.MemStorage{}

			if err := appStorage.Init(); err != nil {
				t.Fatal("Не удалось создать хранилище")
			}

			appStorage.UpdateCounter(tt.args.name, tt.args.counterValue)
			appStorage.UpdateGauge(tt.args.name, tt.args.gaugeValue)

			ch := &UpdateHandler{
				Storage: appStorage,
			}

			request := httptest.NewRequest(tt.args.method, tt.args.addr, nil)
			w := httptest.NewRecorder()

			ch.ServeHTTP(w, request)

			res := w.Result()

			assert.Equal(t, tt.want.code, res.StatusCode)

			defer res.Body.Close()

			if tt.args.name != "" {
				valCounter, _ := ch.Storage.GetCounter(tt.args.name)
				assert.Equal(t, tt.want.counterValue, valCounter)

				valGauge, _ := ch.Storage.GetGauge(tt.args.name)
				assert.Equal(t, tt.want.gaugeValue, valGauge)
			}
		})
	}
}
