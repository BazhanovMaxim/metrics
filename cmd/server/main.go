package main

import (
	"github.com/BazhanovMaxim/metrics/internal/server/configs"
	"github.com/BazhanovMaxim/metrics/internal/server/handlers"
	"github.com/BazhanovMaxim/metrics/internal/server/logger"
	"github.com/BazhanovMaxim/metrics/internal/server/metrics"
	"github.com/BazhanovMaxim/metrics/internal/server/service"
	"github.com/BazhanovMaxim/metrics/internal/server/storage"
	"go.uber.org/zap"
)

func main() {
	config, err := configs.NewConfig()
	if err != nil {
		panic(err)
	}

	if err = logger.NewLogger("info"); err != nil {
		panic(err)
	}

	memStorage := storage.NewMemStorage()
	fileStorage := storage.NewFileStorage(config.FileStoragePath + config.FileStorageName)
	//dbStorage := *storage.NewDBStorage()
	metricService := *service.NewMetricService(config, memStorage, fileStorage)
	if config.Restore {
		metricService.LoadStorageMetrics()
	}

	if config.StoreInterval != 0 {
		go metrics.StartPeriodicSave(config, metricService)
	}

	logger.Log.Info("Running server", zap.String("address", config.RunAddress))
	if handlerError := handlers.NewHandler(config, metricService).Start(); handlerError != nil {
		zap.Error(handlerError)
	}

	// Сохранение всех метрик по завершению
	logger.Log.Info("Saving metrics upon completion server")
	defer metricService.SaveMetricsToStorage()
}
