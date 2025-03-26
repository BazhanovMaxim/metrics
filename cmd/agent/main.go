package main

import (
	"context"
	"github.com/BazhanovMaxim/metrics/internal/agent/configs"
	"github.com/BazhanovMaxim/metrics/internal/agent/logger"
	"github.com/BazhanovMaxim/metrics/internal/agent/router"
	"github.com/BazhanovMaxim/metrics/internal/agent/service"
	"github.com/BazhanovMaxim/metrics/internal/agent/storage"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"time"
)

func main() {
	client := resty.New()
	if err := logger.NewLogger(client); err != nil {
		panic(err)
	}

	config, configError := configs.NewConfig()
	if configError != nil {
		logger.Log.Error("Failed to build config client", zap.Error(configError))
		return
	}

	logger.Log.Info("Running agent", zap.String("address", config.RunAddress))

	// Контекст работы агента
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.AgentWorkingTime)*time.Second)
	defer cancel()

	go func() {
		service.NewMetricService(config, storage.NewMetricRepository(), *router.NewRouter(client)).Start(ctx)
	}()

	// Ждём завершения контекста
	<-ctx.Done()
	logger.Log.Info("The agent has completed the work")
}
