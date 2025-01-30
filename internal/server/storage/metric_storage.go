package storage

// todo: непотокобезопасная
type MetricStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

func NewMetricRepository() IMetricStorage {
	return &MetricStorage{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}
}

type IMetricStorage interface {
	UpdateGauge(key string, value float64)
	UpdateCounter(key string, value int64)
}

func (storage *MetricStorage) UpdateGauge(key string, value float64) {
	storage.Gauge[key] = value
}

func (storage *MetricStorage) UpdateCounter(key string, value int64) {
	if val, find := storage.Counter[key]; find {
		storage.Counter[key] = val + value
		return
	}
	storage.Counter[key] = value
}
