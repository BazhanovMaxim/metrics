package main

import (
	"github.com/BazhanovMaxim/metrics/internal/server/handler"
	"github.com/BazhanovMaxim/metrics/internal/server/repository"
)

func main() {
	err := handler.NewHandler(repository.NewMetricRepository()).Start()
	if err != nil {
		panic(err)
	}
}
