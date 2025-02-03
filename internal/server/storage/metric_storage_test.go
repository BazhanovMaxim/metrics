package storage

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestMetricStorage_UpdateCounter(t *testing.T) {
	tests := []struct {
		name        string
		error       string
		key         string
		counter     int64
		expectedSum int64
	}{
		{"Positive check update counter", "The value has not changed", "first", 1, 1},
		{"Positive check update counter", "The value has not changed", "first", 99, 100},
		{"Positive check update counter", "The value has not changed", "first", -50, 50},
	}
	storage := NewMetricRepository()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage.Counter.Update(test.key, test.counter)
			mapValues := readCounterData(storage)
			assert.Equal(t, test.expectedSum, mapValues[test.key], test.error)
		})
	}
}

func TestMetricStorage_UpdateGauge(t *testing.T) {
	tests := []struct {
		name     string
		error    string
		key      string
		counter  float64
		expected float64
	}{
		{"Positive check update gauge", "The value has not changed", "first", 1, 1},
		{"Positive check update gauge", "The value has not changed", "first", -50, -50},
	}
	storage := NewMetricRepository()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage.Gauge.Update(test.key, test.counter)
			mapValues := readGaugeData(storage)
			assert.Equal(t, test.expected, mapValues[test.key], test.error)
		})
	}
}

func readCounterData(metricStorage *MetricStorage) map[string]int64 {
	return readField("Counter", metricStorage).Interface().(map[string]int64)
}

func readGaugeData(metricStorage *MetricStorage) map[string]float64 {
	return readField("Gauge", metricStorage).Interface().(map[string]float64)
}

func readField(fieldName string, metricStorage *MetricStorage) reflect.Value {
	// Получаем значение поля Counter
	field := reflect.ValueOf(metricStorage).Elem().FieldByName(fieldName)
	// Получаем конкретное значение, на которое указывает интерфейс
	storage := reflect.ValueOf(field.Interface()).Elem()
	// Получаем поле Data
	return storage.FieldByName("Data")
}
