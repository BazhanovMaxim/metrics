package service

import (
	"github.com/BazhanovMaxim/metrics/internal/server/configs"
	"github.com/BazhanovMaxim/metrics/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMetricService_FindService(t *testing.T) {
	tests := []struct {
		name       string
		metricType string
		expectedOk bool
		errorText  string
	}{
		{"Positive get Gauge service", "gauge", true, "Gauge service is not founded"},
		{"Positive get Counter service", "counter", true, "Counter service is not founded"},
		{"Negative unknown service", "unknown", false, "Unknown service is founded"},
	}
	config, _ := configs.NewConfig()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := NewMetricService(config, nil, nil).FindService(test.metricType)
			if test.expectedOk {
				assert.NotNil(t, service)
				return
			}
			assert.Nil(t, service)
		})
	}
}

func TestMetricService_GetMetricValue(t *testing.T) {
	tests := []struct {
		name       string
		metricType string
		expectedOk bool
		errorText  string
	}{
		{"Positive get Gauge", "gauge", true, "Gauge is not founded"},
		{"Positive get Counter", "counter", true, "Counter is not founded"},
		{"Negative no value", "unknown", false, "Unknown service is founded"},
	}
	config, _ := configs.NewConfig()
	memStorage := storage.NewMemStorage()
	memStorage.UpdateGauge("gauge", 10)
	memStorage.UpdateCounter("counter", 10)
	metricService := NewMetricService(config, memStorage, nil)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := metricService.GetMetricValue(test.metricType, test.metricType)
			if test.expectedOk {
				assert.NotNil(t, service)
				return
			}
			assert.Nil(t, service)
		})
	}
}

func TestMetricService_GetCountersGauges(t *testing.T) {
	tests := []struct {
		name          string
		metricType    string
		expectedSize  int
		expectedValue interface{}
		errorText     string
	}{
		{"Positive get Gauges", "gauge", 1, float64(10), "Gauge is not founded"},
		{"Positive get Counters", "counter", 1, int64(10), "Counter is not founded"},
	}

	config, _ := configs.NewConfig()
	memStorage := storage.NewMemStorage()
	memStorage.UpdateGauge("gauge", 10)
	memStorage.UpdateCounter("counter", 10)
	metricService := NewMetricService(config, memStorage, nil)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.metricType == "gauge" {
				service := metricService.GetGauges()
				assert.NotNil(t, service, test.errorText)
				assert.Equal(t, test.expectedValue, service[test.metricType], test.errorText)
				assert.True(t, len(service) == test.expectedSize, test.errorText)
				return
			}
			service := metricService.GetCounters()
			assert.NotNil(t, service, test.errorText)
			assert.Equal(t, test.expectedValue, service[test.metricType], test.errorText)
			assert.True(t, len(service) == test.expectedSize, test.errorText)
		})
	}
}

// todo: переписать на мок
//func TestMetricService_SaveMetricsToStorage(t *testing.T) {
//	type GaugeTest struct {
//		size  int
//		key   string
//		value float64
//	}
//	type CounterTest struct {
//		size  int
//		key   string
//		value int64
//	}
//	tests := []struct {
//		name      string
//		errorText string
//		gauge     GaugeTest
//		counter   CounterTest
//	}{
//		{
//			"Positive SaveMetricsToStorage", "Failed to save metrics to storage",
//			GaugeTest{size: 1, key: "gauge", value: 10},
//			CounterTest{size: 1, key: "counter", value: 10},
//		},
//	}
//
//	filePath := "test.json"
//	file, _ := createTestFile(filePath)
//	writeToFileData(filePath, []byte("{}"))
//	defer file.Close()
//	defer removeTestFile(filePath)
//
//	metricRepository := storage.NewMemStorage()
//	metricRepository.UpdateGauge("gauge", 10)
//	metricRepository.UpdateCounter("counter", 10)
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			NewMetricService().SaveMetricsToStorage()
//			metrics := model.StorageJSONMetrics{}
//			data, _ := os.ReadFile(filePath)
//			_ = json.Unmarshal(data, &metrics)
//			// gauge
//			assert.Equal(t, test.gauge.size, len(metrics.Gauge), test.errorText)
//			assert.Equal(t, test.gauge.value, metrics.Gauge[test.gauge.key], test.errorText)
//			//counter
//			assert.Equal(t, test.counter.size, len(metrics.Counter), test.errorText)
//			assert.Equal(t, test.counter.value, metrics.Counter[test.counter.key], test.errorText)
//		})
//	}
//}

// todo: переписать на мок
//func TestMetricService_SaveMetricToStorage(t *testing.T) {
//	tests := []struct {
//		name      string
//		errorText string
//		metrics   []model.Metrics
//	}{
//		{
//			name:      "Positive SaveMetricsToStorage",
//			errorText: "Failed to save metrics to storage",
//			metrics: []model.Metrics{
//				{ID: "gauge", MType: "gauge", Value: float64Pointer(20)},
//				{ID: "counter", MType: "counter", Delta: int64Pointer(20)},
//			},
//		},
//	}
//
//	filePath := "test.json"
//	file, _ := createTestFile(filePath)
//	writeToFileData(filePath, []byte("{}"))
//	defer file.Close()
//	defer removeTestFile(filePath)
//
//	for _, test := range tests {
//		for _, mtr := range test.metrics {
//			t.Run(test.name, func(t *testing.T) {
//				NewMetricService().SaveMetricToStorage(filePath, &mtr)
//				metrics := model.StorageJSONMetrics{}
//				data, _ := os.ReadFile(filePath)
//				_ = json.Unmarshal(data, &metrics)
//				if mtr.MType == "counter" {
//					assert.Equal(t, *mtr.Delta, metrics.Counter[mtr.ID], test.errorText)
//					return
//				}
//				assert.Equal(t, *mtr.Value, metrics.Gauge[mtr.ID], test.errorText)
//			})
//		}
//	}
//}

// Вспомогательная функция для создания указателя на int64
func int64Pointer(i int) *int64 {
	value := int64(i)
	return &value
}

func float64Pointer(i float64) *float64 {
	return &i
}

func createTestFile(filePath string) (*os.File, error) {
	return os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
}

func writeToFileData(filePath string, data []byte) {
	_ = os.WriteFile(filePath, data, 0666)
}

func removeTestFile(filePath string) {
	_ = os.Remove(filePath)
}
