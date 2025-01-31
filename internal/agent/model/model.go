package model

import "fmt"

type Metric struct {
	Value interface{}
	Type  MetricType
}

type MetricType string

func (metric *Metric) ValueAsString() string {
	switch metric.Type {
	case Gauge:
		// gauge
		return fmt.Sprintf("%f", metric.Value.(float64))
	default:
		// counter
		return fmt.Sprintf("%d", metric.Value.(int64))
	}
}

const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)
