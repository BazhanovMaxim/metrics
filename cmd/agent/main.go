package main

import (
	"github.com/BazhanovMaxim/metrics/internal/agent/flags"
	"github.com/BazhanovMaxim/metrics/internal/agent/handlers"
	"github.com/BazhanovMaxim/metrics/internal/agent/storage"
)

func main() {
	flags.ParseAgentFlags()
	err := handlers.NewHandler(storage.NewMetricRepository()).Start()
	if err != nil {
		panic(err)
	}
}
