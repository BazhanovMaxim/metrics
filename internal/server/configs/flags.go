package configs

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

var FlagRunAddr string

type Config struct {
	RunAddr string `env:"ADDRESS"`
}

func ParseServerConfigs() error {
	// 1. Парсим флаги
	parseServerFlags()
	// 2. Парсим переменные окружения
	if err := parseOsEnv(); err != nil {
		return err
	}
	return nil
}

func parseServerFlags() {
	flag.StringVar(&FlagRunAddr, "a", ":8080", "address and port to run server")
	flag.Parse()
}

func parseOsEnv() error {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return err
	}
	if cfg.RunAddr != "" {
		FlagRunAddr = cfg.RunAddr
	}
	return nil
}
