package storage

type IMetricStorage[metricType float64 | int64] interface {
	Update(key string, value metricType)
	Get(key string) (metricType, bool)
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

func (s *CounterStorage) Update(key string, value int64) {
	if val, find := s.Data[key]; find {
		s.Data[key] = val + value
		return
	}
	s.Data[key] = value
}

func (s *CounterStorage) Get(key string) (int64, bool) {
	if val, find := s.Data[key]; find {
		return val, find
	}
	return 0, false
}

func (s *CounterStorage) GetAll() map[string]int64 {
	return s.Data
}

func (s *GaugeStorage) Update(key string, value float64) {
	s.Data[key] = value
}

func (s *GaugeStorage) Get(key string) (float64, bool) {
	if val, find := s.Data[key]; find {
		return val, find
	}
	return 0.0, false
}

func (s *GaugeStorage) GetAll() map[string]float64 {
	return s.Data
}
