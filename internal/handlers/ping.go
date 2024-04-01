package handlers

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/KirillKhitev/metrics/internal/logger"
)

// PingHandler отвечает за обработку GET-запроса /ping.
type PingHandler MyHandler

/*
ServeHTTP служит для проверки соединения с БД.

Коды ответа:

• 200 - успешный запрос.

• 405 - метод запроса отличен от GET.

• 500 - при возникновении ошибки.
*/
func (ch *PingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err := ch.Storage.Ping(r.Context()); err != nil {
		logger.Log.Error("Cannot connect to db", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
