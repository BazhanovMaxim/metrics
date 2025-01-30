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
	GetGauge(key string) (float64, bool)
	GetCounter(key string) (int64, bool)
	GetGauges() map[string]float64
	GetCounters() map[string]int64
}

func (storage *MetricStorage) UpdateGauge(key string, value float64) {
	storage.Gauge[key] = value
}

func (storage *MetricStorage) GetGauge(key string) (float64, bool) {
	if val, find := storage.Gauge[key]; find {
		return val, true
	}
	return 0.0, false
}

func (storage *MetricStorage) GetGauges() map[string]float64 {
	return storage.Gauge
}

func (storage *MetricStorage) UpdateCounter(key string, value int64) {
	if val, find := storage.Counter[key]; find {
		storage.Counter[key] = val + value
		return
	}
	storage.Counter[key] = value
}

func (storage *MetricStorage) GetCounter(key string) (int64, bool) {
	if val, find := storage.Counter[key]; find {
		return val, true
	}
	return 0.0, false
}

func (storage *MetricStorage) GetCounters() map[string]int64 {
	return storage.Counter
}
