package repository

// todo: непотокобезопасная
type MetricStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func NewMetricRepository() *MetricStorage {
	return &MetricStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

type iMetricStorage interface {
	UpdateGauge(key string, value float64)
	UpdateCounter(key string, value int64)
}

func (storage *MetricStorage) UpdateGauge(key string, value float64) {
	storage.gauge[key] = value
}

func (storage *MetricStorage) UpdateCounter(key string, value int64) {
	if val, find := storage.counter[key]; find {
		storage.counter[key] = val + value
		return
	}
	storage.counter[key] = value
}
