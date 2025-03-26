package storage

import (
	"github.com/BazhanovMaxim/metrics/internal/server/model"
	"github.com/BazhanovMaxim/metrics/internal/server/service"
	"sync"
)

// MemStorage представляет собой локальное хранилище для работы с метриками
type MemStorage struct {
	mutex   sync.Mutex
	Counter map[string]int64
	Gauge   map[string]float64
}

// NewMemStorage создает и возвращает новый экземпляр MemStorage
func NewMemStorage() service.IMetricStorage {
	return &MemStorage{
		Counter: make(map[string]int64),
		Gauge:   make(map[string]float64),
	}
}

// Update обновляет метрики в локальном хранилище. Если метрики были добавлены ранее,
// значение этих метрик будут изменены, иначе будет добавлена новая запись метрики
func (m *MemStorage) Update(metric model.Metrics) (*model.Metrics, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.updateMetric(metric)
}

// UpdateBatches обновляет метрики в локальном хранилище. Если метрики были добавлены ранее,
// значение этих метрик будут изменены, иначе будет добавлена новая запись метрики
func (m *MemStorage) UpdateBatches(metrics []model.Metrics) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for _, metric := range metrics {
		if _, err := m.updateMetric(metric); err != nil {
			return err
		}
	}
	return nil
}

func (m *MemStorage) updateMetric(metric model.Metrics) (*model.Metrics, error) {
	if metric.MType == string(model.Gauge) {
		m.Gauge[metric.ID] = *metric.Value
		return &metric, nil
	}
	if val, find := m.Counter[metric.ID]; find {
		newValue := val + *metric.Delta
		m.Counter[metric.ID] = newValue
		return &model.Metrics{ID: metric.ID, MType: metric.MType, Delta: &newValue}, nil
	}
	m.Counter[metric.ID] = *metric.Delta
	return &metric, nil
}

// GetAllMetrics получает и возвращает все метрики из локального хранилища
func (m *MemStorage) GetAllMetrics() []model.Metrics {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var t []model.Metrics
	if len(m.Counter) != 0 {
		for key, value := range m.Counter {
			t = append(t, model.Metrics{ID: key, MType: string(model.Counter), Delta: &value})
		}
	}
	if len(m.Gauge) != 0 {
		for key, value := range m.Gauge {
			t = append(t, model.Metrics{ID: key, MType: string(model.Gauge), Value: &value})
		}
	}
	return t
}

// GetMetric получает и возвращает метрику из локального хранилища.
// В случае, если метрики нет, тогда возвращается nil
func (m *MemStorage) GetMetric(mType, title string) *model.Metrics {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if mType == string(model.Counter) {
		if val, find := m.Counter[title]; find {
			return &model.Metrics{ID: title, MType: mType, Delta: &val}
		}
		return nil
	}
	if val, find := m.Gauge[title]; find {
		return &model.Metrics{ID: title, MType: mType, Value: &val}
	}
	return nil
}

func (m *MemStorage) Ping() error {
	return nil
}

func (m *MemStorage) Close() error {
	return nil
}
