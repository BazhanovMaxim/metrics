package storage

// IMemStorage todo: проблема интерфейса в количестве методов.
// Пробовал с использованием дженериков, но дженерики у интерфейса (IMemStorage[V int64 | float64]) - является
// антипаттерном. Подумать, как можно упростить не теряя в удобстве использования и с возможностью использования
// полиморфизма.
type IMemStorage interface {
	IMetricStorage
	UpdateCounter(key string, value int64) int64
	UpdateGauge(key string, value float64) float64
	GetCounters() map[string]int64
	GetGauges() map[string]float64
}

type MemStorage struct {
	Counter map[string]int64
	Gauge   map[string]float64
}

func NewMemStorage() IMemStorage {
	memStorage := &MemStorage{}
	memStorage.init()
	return memStorage
}

func (m *MemStorage) init() {
	m.Counter = make(map[string]int64)
	m.Gauge = make(map[string]float64)
}

func (m *MemStorage) UpdateCounter(key string, value int64) int64 {
	if val, find := m.Counter[key]; find {
		m.Counter[key] = val + value
		return m.Counter[key]
	}
	m.Counter[key] = value
	return value
}

func (m *MemStorage) UpdateGauge(key string, value float64) float64 {
	m.Gauge[key] = value
	return value
}

func (m *MemStorage) GetCounters() map[string]int64 {
	return m.Counter
}

func (m *MemStorage) GetGauges() map[string]float64 {
	return m.Gauge
}
