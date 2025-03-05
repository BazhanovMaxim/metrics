package storage

import (
	"database/sql"
	"github.com/BazhanovMaxim/metrics/internal/server/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

type IDbStorage interface {
	IMetricStorage
	Ping() error
	Close() error
	// Request
}

type DBStorage struct {
	dataSourceName string
	db             *sql.DB
}

func NewDBStorage(dataSourceName string) IDbStorage {
	dbStorage := DBStorage{dataSourceName: dataSourceName}
	dbStorage.init()
	return &dbStorage
}

func (s *DBStorage) init() {
	db, err := sql.Open("pgx", s.dataSourceName)
	if err != nil {
		logger.Log.Error("Failed to connect Database with raw connect", zap.Error(err))
		return
	}
	s.db = db
	if pingError := s.Ping(); pingError != nil {
		logger.Log.Error("Failed to connect Database with raw connect", zap.Error(pingError))
		return
	}
	logger.Log.Info("Success database connection")
}

func (s *DBStorage) Ping() error {
	return s.db.Ping()
}

func (s *DBStorage) Close() error {
	return s.db.Close()
}
