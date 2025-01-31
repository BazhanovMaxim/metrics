package configs

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"os"
)

var (
	RunAddress     string
	ReportInterval int
	PollInterval   int
)

type Config struct {
	RunAddr string `env:"ADDRESS"`
	Report  int    `env:"REPORT_INTERVAL"`
	Poll    int    `env:"POLL_INTERVAL"`
}

func ParseAgentConfigs() error {
	// 1. Парсим флаги
	parseAgentFlags()
	// 2. Парсим переменные окружения
	if err := parseOsEnv(); err != nil {
		return err
	}
	return nil
}

func parseAgentFlags() {
	flagSet := flag.NewFlagSet("agentFlags", flag.ContinueOnError)
	flagSet.StringVar(&RunAddress, "a", "localhost:8080", "address and port to run agent")
	flagSet.IntVar(&ReportInterval, "r", 10, "report interval")
	flagSet.IntVar(&PollInterval, "p", 2, "poll interval")
	flagSet.Parse(os.Args[1:])
}

func parseOsEnv() error {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return err
	}
	if cfg.RunAddr != "" {
		RunAddress = cfg.RunAddr
	}
	if cfg.Poll != 0 {
		PollInterval = cfg.Poll
	}
	if cfg.Report != 0 {
		ReportInterval = cfg.Report
	}
	return nil
}
