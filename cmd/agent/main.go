package main

import (
	"github.com/BazhanovMaxim/metrics/internal/agent/configs"
	"github.com/BazhanovMaxim/metrics/internal/agent/handlers"
	"github.com/BazhanovMaxim/metrics/internal/agent/storage"
)

func main() {
	config, configError := configs.NewConfig()
	if configError != nil {
		panic(configError)
	}
	err := handlers.NewHandler(config, storage.NewMetricRepository()).Start()
	if err != nil {
		panic(err)
	}
}
