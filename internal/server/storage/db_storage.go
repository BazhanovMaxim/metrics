package storage

import (
	"context"
	"errors"
	"github.com/BazhanovMaxim/metrics/internal/server/configs"
	"github.com/BazhanovMaxim/metrics/internal/server/logger"
	"github.com/BazhanovMaxim/metrics/internal/server/model"
	"github.com/BazhanovMaxim/metrics/internal/server/queries"
	"github.com/BazhanovMaxim/metrics/internal/server/service"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

// DBStorage представляет собой хранилище для работы с метриками в базе данных
type DBStorage struct {
	config configs.Config
	db     *sqlx.DB
}

// NewDBStorage создает и возвращает новый экземпляр DBStorage
func NewDBStorage(config configs.Config) service.IMetricStorage {
	dbStorage := DBStorage{config: config}
	dbStorage.init()
	return &dbStorage
}

// init инициализирует контекст DBStorage
func (s *DBStorage) init() {
	// открываем пул соединений и сразу проверяем работоспособность подключения
	db, err := sqlx.Connect("pgx", s.config.DatabaseDSN)
	if err != nil {
		logger.Log.Info("Database connection error. Retry to reconnect to database")
		// выполняет попытку переподключения к базе данных
		ticker := time.NewTicker(2 * time.Second)
		for i := 0; i < 3; i++ {
			<-ticker.C
			if db, err = sqlx.Connect("pgx", s.config.DatabaseDSN); err == nil {
				s.db = db
				ticker.Stop()
				break
			}
		}
		ticker.Stop()
		logger.Log.Error("Database connection error")
		return
	}
	s.db = db
	// Настройка параметров пула соединений
	s.db.SetMaxOpenConns(s.config.MaxOpenCons)                                     // Максимальное количество открытых соединений
	s.db.SetMaxIdleConns(s.config.MaxIdleCons)                                     // Максимальное количество простаивающих соединений
	s.db.SetConnMaxLifetime(time.Duration(s.config.ConnMaxLifetime) * time.Minute) // Максимальное время жизни соединения
	s.db.SetConnMaxIdleTime(time.Duration(s.config.ConnMaxIdleTime) * time.Minute) // Максимальное время простоя соединения
	s.createMetricsTable()
	logger.Log.Info("Success database connection")
}

// Update обновляет метрики в базе данных. Если метрики были добавлены ранее,
// значение этих метрик будут изменены, иначе будет добавлена новая запись метрики
func (s *DBStorage) Update(metric model.Metrics) (*model.Metrics, error) {
	_, err := s.db.Exec(queries.InsertMetric, metric.ID, metric.MType, metric.Delta, metric.Value)
	if err != nil {
		logger.Log.Error("Failed to update metric", zap.Error(err))
		return nil, err
	}
	return &metric, err
}

// UpdateBatches выполняет транзакцию для множественной вставки метрик в базе данных. Если метрики были добавлены ранее,
// значение этих метрик будут изменены, иначе будут добавлены новые записи метрик
func (s *DBStorage) UpdateBatches(metrics []model.Metrics) error {
	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		logger.Log.Error("Failed to create transaction", zap.Error(err))
		return err
	}
	for _, metric := range metrics {
		_, insErr := tx.Exec(queries.InsertMetric, metric.ID, metric.MType, metric.Delta, metric.Value)
		if insErr != nil {
			logger.Log.Error("Failed to update batches", zap.Error(insErr))
			_ = tx.Rollback()
			return insErr
		}
	}
	return tx.Commit()
}

// GetAllMetrics получает и возвращает данных о всех метриках из базы данных
func (s *DBStorage) GetAllMetrics() []model.Metrics {
	var metric []model.Metrics
	if err := s.db.Select(&metric, queries.GetMetrics); err != nil {
		logger.Log.Error("Failed get all metrics by database", zap.Error(err))
		return nil
	}
	return metric
}

// GetMetric получает и возвращает данные о какой-то определенной метрики из базы данных.
// Возвращает nil в случае, если искомой метрики нет в базе данных
func (s *DBStorage) GetMetric(mType, title string) *model.Metrics {
	var metric model.Metrics
	if err := s.db.Get(&metric, queries.GetMetric, title, mType); err != nil {
		logger.Log.Error("Failed get metric by database", zap.Error(err))
		return nil
	}
	return &metric
}

// Close закрывает соединение с базой данных
func (s *DBStorage) Close() error {
	if s.db == nil {
		return errors.New("database connection is empty")
	}
	return s.db.Close()
}

// createMetricsTable создает таблицу для метрик в базе данных
func (s *DBStorage) createMetricsTable() {
	if _, err := s.db.Exec(queries.CreateMetricsTable); err != nil {
		logger.Log.Error("Failed to create metrics table in database", zap.Error(err))
		return
	}
	if _, err := s.db.Exec(queries.CreateIndex); err != nil {
		logger.Log.Error("Failed to create metrics index in database", zap.Error(err))
		return
	}
	logger.Log.Info("The metric tables were created successfully")
}

// Ping проверяет работоспособность подключения к базе данных
func (s *DBStorage) Ping() error {
	if s.db == nil {
		return errors.New("database connection is empty")
	}
	return s.db.Ping()
}
