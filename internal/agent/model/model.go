package model

type Metric struct {
	Value interface{}
	Type  MetricType
}

type MetricType string

const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)
