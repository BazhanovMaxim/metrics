package main

import (
	"github.com/BazhanovMaxim/metrics/internal/agent/handlers"
	"github.com/BazhanovMaxim/metrics/internal/agent/storage"
)

func main() {
	err := handlers.NewHandler(storage.NewMetricRepository()).Start()
	if err != nil {
		panic(err)
	}
}
