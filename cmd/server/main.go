package main

import (
	"github.com/BazhanovMaxim/metrics/internal/server/configs"
	"github.com/BazhanovMaxim/metrics/internal/server/handlers"
	"github.com/BazhanovMaxim/metrics/internal/server/logger"
	"github.com/BazhanovMaxim/metrics/internal/server/storage"
	"go.uber.org/zap"
)

func main() {
	config, err := configs.NewConfig()
	if err != nil {
		panic(err)
	}

	if err := logger.NewLogger("info"); err != nil {
		panic(err)
	}

	logger.Log.Info("Running server", zap.String("address", config.RunAddress))
	if err := handlers.NewHandler(config, *storage.NewMetricRepository()).Start(); err != nil {
		zap.Error(err)
	}
}
