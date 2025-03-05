package configs

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"os"
)

type OsConfig struct {
	RunAddr         string `env:"ADDRESS"`
	StoreInterval   int64  `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	FileName        string `env:"FILE_NAME"`
	Restore         bool   `env:"RESTORE"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
}

type Config struct {
	RunAddress string
	// Интервал времени в секундах, по истечении которого текущие показания
	// сервера сохраняются на диск (по умолчанию 300 секунд, значение 0 делает запись синхронной)
	StoreInterval int64
	// Путь до файла, куда сохраняются текущие значения. Имя файла для значения по умолчанию придумайте сами
	FileStoragePath string
	// Имя файла
	FileStorageName string
	// Загружать или нет ранее сохранённые значения из указанного файла при старте сервера
	Restore bool
	// Строка с адресом подключения к БД
	DatabaseDSN string
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
	flagSet.Int64Var(&config.StoreInterval, "i", 300, "store interval")
	flagSet.StringVar(&config.FileStoragePath, "f", "internal/server/tmp/test.json", "file storage path")
	flagSet.StringVar(&config.FileStorageName, "n", "/test.json", "file name")
	flagSet.BoolVar(&config.Restore, "r", false, "load saved metric value when the server starts")
	flagSet.StringVar(&config.DatabaseDSN, "d", "", "database URL connection")
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
	if cfg.StoreInterval != 0 {
		config.StoreInterval = cfg.StoreInterval
	}
	if cfg.Restore {
		config.Restore = cfg.Restore
	}
	if cfg.FileStoragePath != "" {
		config.FileStoragePath = cfg.FileStoragePath
	}
	if cfg.FileName != "" {
		config.FileStorageName = cfg.FileName
	}
	if cfg.DatabaseDSN != "" {
		config.DatabaseDSN = cfg.DatabaseDSN
	}
	return nil
}
