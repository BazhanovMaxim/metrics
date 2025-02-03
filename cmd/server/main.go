package main

import (
	"github.com/BazhanovMaxim/metrics/internal/server/configs"
	"github.com/BazhanovMaxim/metrics/internal/server/handlers"
	"github.com/BazhanovMaxim/metrics/internal/server/storage"
	"log"
)

func main() {
	config, err := configs.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	if err := handlers.NewHandler(config, *storage.NewMetricRepository()).Start(); err != nil {
		log.Fatal(err)
	}
}
