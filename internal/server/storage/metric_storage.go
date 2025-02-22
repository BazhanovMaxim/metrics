package storage

type IMetricStorage[metricType float64 | int64] interface {
	Update(key string, value metricType) metricType
	GetAll() map[string]metricType
}

type MetricStorage struct {
	Counter IMetricStorage[int64]
	Gauge   IMetricStorage[float64]
}

func NewMetricRepository() *MetricStorage {
	return &MetricStorage{
		Counter: &CounterStorage{Data: make(map[string]int64)},
		Gauge:   &GaugeStorage{Data: make(map[string]float64)},
	}
}

type CounterStorage struct {
	Data map[string]int64
}

type GaugeStorage struct {
	Data map[string]float64
}

func (s *CounterStorage) Update(key string, value int64) int64 {
	if val, find := s.Data[key]; find {
		s.Data[key] = val + value
		return s.Data[key]
	}
	s.Data[key] = value
	return value
}

func (s *CounterStorage) GetAll() map[string]int64 {
	return s.Data
}

func (s *GaugeStorage) Update(key string, value float64) float64 {
	s.Data[key] = value
	return value
}

func (s *GaugeStorage) GetAll() map[string]float64 {
	return s.Data
}
