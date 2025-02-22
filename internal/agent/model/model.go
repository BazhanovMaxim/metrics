package model

// Metrics updateMetric metrics model
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type Metric struct {
	Value interface{}
	Type  MetricType
}

type MetricType string

const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)
