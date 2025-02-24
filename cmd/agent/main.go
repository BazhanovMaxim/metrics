package main

import (
	"github.com/BazhanovMaxim/metrics/internal/agent/configs"
	"github.com/BazhanovMaxim/metrics/internal/agent/handlers"
	"github.com/BazhanovMaxim/metrics/internal/agent/logger"
	"github.com/BazhanovMaxim/metrics/internal/agent/service"
	"github.com/BazhanovMaxim/metrics/internal/agent/storage"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

func main() {
	client := resty.New()
	if err := logger.NewLogger(client); err != nil {
		panic(err)
	}

	config, configError := configs.NewConfig()
	if configError != nil {
		logger.Log.Error("Failed to build client", zap.Error(configError))
	}

	logger.Log.Info("Running agent", zap.String("address", config.RunAddress))
	service.NewMetricService(config, storage.NewMetricRepository(), *handlers.NewHandler(client)).Start()
}
