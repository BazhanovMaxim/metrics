package configs

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"os"
)

type OsConfig struct {
	RunAddr string `env:"ADDRESS"`
}

type Config struct {
	RunAddress string
}

func NewConfig() (Config, error) {
	config := Config{}
	// 1. Парсим флаги
	parseServerFlags(&config)
	// 2. Парсим переменные окружения
	if err := parseOsEnv(&config); err != nil {
		return config, err
	}
	return config, nil
}

func parseServerFlags(config *Config) {
	flagSet := flag.NewFlagSet("serverFlags", flag.ContinueOnError)
	flagSet.StringVar(&config.RunAddress, "a", ":8080", "address and port to run server")
	flagSet.Parse(os.Args[1:])
}

func parseOsEnv(config *Config) error {
	var cfg OsConfig
	if err := env.Parse(&cfg); err != nil {
		return err
	}
	if cfg.RunAddr != "" {
		config.RunAddress = cfg.RunAddr
	}
	return nil
}
