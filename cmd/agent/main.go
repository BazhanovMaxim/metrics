package main

import (
	"github.com/BazhanovMaxim/metrics/internal/agent/configs"
	"github.com/BazhanovMaxim/metrics/internal/agent/handlers"
	"github.com/BazhanovMaxim/metrics/internal/agent/storage"
)

func main() {
	if err := configs.ParseAgentConfigs(); err != nil {
		panic(err)
	}
	err := handlers.NewHandler(storage.NewMetricRepository()).Start()
	if err != nil {
		panic(err)
	}
}
