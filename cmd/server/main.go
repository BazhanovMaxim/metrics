package main

import (
	"github.com/BazhanovMaxim/metrics/internal/server/configs"
	"github.com/BazhanovMaxim/metrics/internal/server/handlers"
	"github.com/BazhanovMaxim/metrics/internal/server/storage"
	"log"
)

func main() {
	if err := configs.ParseServerConfigs(); err != nil {
		log.Fatal(err)
	}
	if err := handlers.NewHandler(storage.NewMetricRepository()).Start(); err != nil {
		log.Fatal(err)
	}
}
