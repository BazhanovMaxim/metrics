package main

import (
	"github.com/BazhanovMaxim/metrics/internal/server/configs"
	"github.com/BazhanovMaxim/metrics/internal/server/handlers"
	"github.com/BazhanovMaxim/metrics/internal/server/logger"
	"github.com/BazhanovMaxim/metrics/internal/server/service"
	"github.com/BazhanovMaxim/metrics/internal/server/storage"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	config, err := configs.NewConfig()
	if err != nil {
		panic(err)
	}

	if err = logger.NewLogger("info"); err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	// Создаем канал для получения сигналов
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	rep := newStorage(config)
	closeApp(&wg, sigChan, rep)

	metricService := *service.NewMetricService(config, rep)

	logger.Log.Info("Running server", zap.String("address", config.RunAddress))
	if handlerError := handlers.NewHandler(config, metricService).Start(); handlerError != nil {
		zap.Error(handlerError)
	}

	wg.Wait()
	logger.Log.Info("Server has been stopped")
}

func closeApp(wg *sync.WaitGroup, sigChan chan os.Signal, storage service.IMetricStorage) {
	wg.Add(1)

	// Запускаем горутину, которая будет выполнять действие при завершении программы
	go func() {
		defer wg.Done()
		// Ждем сигнала завершения
		<-sigChan
		_ = storage.Close()
		os.Exit(0)
	}()
}

func newStorage(config configs.Config) service.IMetricStorage {
	if config.DatabaseDSN != "" {
		logger.Log.Info("The database source is being used with url: " + config.DatabaseDSN)
		return storage.NewDBStorage(config.DatabaseDSN)
	}
	if config.FileStoragePath != "" {
		logger.Log.Info("The file source is being used with file path: " + config.FileStoragePath + config.FileStorageName)
		return storage.NewFileStorage(config)
	}
	logger.Log.Info("The memory source is being used")
	return storage.NewMemStorage()
}
