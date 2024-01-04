package dump

import (
	"encoding/json"
	"github.com/KirillKhitev/metrics/internal/logger"
	"github.com/KirillKhitev/metrics/internal/storage"
	"go.uber.org/zap"
	"os"
)

func SaveStorageToFile(filepath string, appStorage storage.MemStorage) {
	logger.Log.Info("Сохраняем метрики в файл")

	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Log.Error("Error by export metrics to file from server", zap.Error(err))
		return
	}

	defer file.Close()

	if err := json.NewEncoder(file).Encode(appStorage); err != nil {
		logger.Log.Error("Error by encode metrics to json", zap.Error(err))
		return
	}
}
