package storage

import (
	"github.com/BazhanovMaxim/metrics/internal/server/logger"
	"github.com/BazhanovMaxim/metrics/internal/server/model"
	"github.com/BazhanovMaxim/metrics/internal/server/queries"
	"github.com/BazhanovMaxim/metrics/internal/server/service"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type DBStorage struct {
	dataSourceName string
	db             *sqlx.DB
}

func (s *DBStorage) Update(metric model.Metrics) (*model.Metrics, error) {
	_, err := s.db.Exec(queries.InsertMetric, metric.ID, metric.MType, metric.Delta, metric.Value)
	if err != nil {
		logger.Log.Error("Failed to update metric", zap.Error(err))
		return nil, err
	}
	return &metric, err
}

func (s *DBStorage) GetAllMetrics() []model.Metrics {
	var metric []model.Metrics
	if err := s.db.Select(metric, queries.GetMetrics); err != nil {
		logger.Log.Error("Failed get metric by database", zap.Error(err))
		return nil
	}
	return metric
}

func (s *DBStorage) GetMetric(mType, title string) *model.Metrics {
	var metric model.Metrics
	if err := s.db.Get(metric, queries.GetMetric, title, mType); err != nil {
		logger.Log.Error("Failed get metric by database", zap.Error(err))
		return nil
	}
	return &metric
}

func (s *DBStorage) Close() error {
	return s.db.Close()
}

func NewDBStorage(dataSourceName string) service.IMetricStorage {
	dbStorage := DBStorage{dataSourceName: dataSourceName}
	dbStorage.init()
	return &dbStorage
}

func (s *DBStorage) init() {
	db, err := sqlx.Open("pgx", s.dataSourceName)
	if err != nil {
		logger.Log.Error("Database connection error", zap.Error(err))
		return
	}
	s.db = db
	// проверяем работоспособность подключения к базе данных
	if pingError := s.Ping(); pingError != nil {
		logger.Log.Error("Database connection error", zap.Error(pingError))
		return
	}
	s.createMetricsTable()
	logger.Log.Info("Success database connection")
}

func (s *DBStorage) createMetricsTable() {
	if _, err := s.db.Exec(queries.CreateMetricsTable); err != nil {
		logger.Log.Error("Failed to create metrics table in database")
		return
	}
	logger.Log.Info("The metric tables were created successfully")
}

func (s *DBStorage) Ping() error {
	return s.db.Ping()
}
