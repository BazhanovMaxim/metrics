package storage

import (
	"github.com/BazhanovMaxim/metrics/internal/server/model"
	"github.com/BazhanovMaxim/metrics/internal/server/service"
)

type MemStorage struct {
	Counter map[string]int64
	Gauge   map[string]float64
}

func NewMemStorage() service.IMetricStorage {
	return &MemStorage{
		Counter: make(map[string]int64),
		Gauge:   make(map[string]float64),
	}
}

func (m *MemStorage) Update(metric model.Metrics) (*model.Metrics, error) {
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

func (m *MemStorage) GetAllMetrics() []model.Metrics {
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

func (m *MemStorage) GetMetric(mType, title string) *model.Metrics {
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
