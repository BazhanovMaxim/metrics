package storage

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestNewMetricRepository(t *testing.T) {
	tests := []struct {
		name         string
		error        string
		fieldName    string
		expectedSize int
		mapType      reflect.Type
	}{
		{"Positive check Counter map is empty", "Counter map is not empty", "Counter", 0, reflect.TypeOf(map[string]int64{})},
		{"Positive check Gauge map is empty", "Gauge map is not empty", "Gauge", 0, reflect.TypeOf(map[string]float64{})},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mapValues := readPrivateMapValues(reflect.ValueOf(NewMetricRepository()).Elem().FieldByName(test.fieldName))
			assert.Equal(t, test.expectedSize, len(mapValues), test.error)
		})
	}
}

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
			storage.UpdateCounter(test.key, test.counter)
			mapValues := readPrivateMapValues(reflect.ValueOf(storage).Elem().FieldByName("Counter"))
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
			storage.UpdateGauge(test.key, test.counter)
			mapValues := readPrivateMapValues(reflect.ValueOf(storage).Elem().FieldByName("Gauge"))
			assert.Equal(t, test.expected, mapValues[test.key], test.error)
		})
	}
}

func readPrivateMapValues(field reflect.Value) map[interface{}]interface{} {
	mapValues := make(map[interface{}]interface{})
	for _, key := range field.MapKeys() {
		value := field.MapIndex(key)
		mapValues[key.Interface()] = value.Interface()
	}
	return mapValues
}
