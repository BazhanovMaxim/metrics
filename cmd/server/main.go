package main

import (
	"github.com/BazhanovMaxim/metrics/internal/server/flags"
	"github.com/BazhanovMaxim/metrics/internal/server/handlers"
	"github.com/BazhanovMaxim/metrics/internal/server/storage"
	"log"
)

func main() {
	flags.ParseFlags()
	err := handlers.NewHandler(storage.NewMetricRepository()).Start()
	if err != nil {
		log.Fatal(err)
	}
}
