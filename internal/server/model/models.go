package model

// IndexHTMLModel Model to generate HTML-page
type IndexHTMLModel struct {
	Key   string
	Value interface{}
}

// Metrics updateMetric metrics model
type Metrics struct {
	ID    string   `db:"mid"   json:"id"`              // Metric name
	MType string   `db:"mtype" json:"type"`            // Gauge or Counter
	Delta *int64   `db:"delta" json:"delta,omitempty"` // Counter value
	Value *float64 `db:"value" json:"value,omitempty"` // Gauge value
}

// StorageJSONMetrics struct for save metrics to file
type StorageJSONMetrics struct {
	Gauge   map[string]float64 `json:"gauge,omitempty"`   // Gauge metrics
	Counter map[string]int64   `json:"counter,omitempty"` // Counter metrics
}

// MetricType constants
type MetricType string

const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)
