package configs

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"os"
)

type Config struct {
	RunAddress       string
	ReportInterval   int
	PollInterval     int
	AgentWorkingTime int
	SecretKey        string
}

type OsConfig struct {
	RunAddr   string `env:"ADDRESS"`
	Report    int    `env:"REPORT_INTERVAL"`
	Poll      int    `env:"POLL_INTERVAL"`
	AgentTime int    `env:"AGENT_WORKING_TIME"`
	Key       string `env:"KEY"`
}

func NewConfig() (Config, error) {
	agentConfig := Config{}
	// 1. Парсим флаги
	parseAgentFlags(&agentConfig)
	// 2. Парсим переменные окружения
	if err := parseOsEnv(&agentConfig); err != nil {
		return agentConfig, err
	}
	return agentConfig, nil
}

func parseAgentFlags(config *Config) {
	flagSet := flag.NewFlagSet("agentFlags", flag.ContinueOnError)
	flagSet.StringVar(&config.RunAddress, "a", "localhost:8080", "address and port to run agent")
	flagSet.IntVar(&config.ReportInterval, "r", 10, "report interval")
	flagSet.IntVar(&config.PollInterval, "p", 2, "poll interval")
	flagSet.IntVar(&config.AgentWorkingTime, "t", 600, "agent's working time")
	flagSet.StringVar(&config.SecretKey, "k", "", "secret key")
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
	if cfg.Poll != 0 {
		config.PollInterval = cfg.Poll
	}
	if cfg.Report != 0 {
		config.ReportInterval = cfg.Report
	}
	if cfg.AgentTime != 0 {
		config.AgentWorkingTime = cfg.AgentTime
	}
	if cfg.Key != "" {
		config.SecretKey = cfg.Key
	}
	return nil
}
