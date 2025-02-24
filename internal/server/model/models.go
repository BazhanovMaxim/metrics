package model

// IndexHTMLModel Model to generate HTML-page
type IndexHTMLModel struct {
	Key   string
	Value interface{}
}

// Metrics updateMetric metrics model
type Metrics struct {
	ID    string   `json:"id"`              // Metric name
	MType string   `json:"type"`            // Gauge or Counter
	Delta *int64   `json:"delta,omitempty"` // Counter value
	Value *float64 `json:"value,omitempty"` // Gauge value
}

// MetricType constants
type MetricType string

const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)
