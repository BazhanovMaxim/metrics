package storage

import (
	"github.com/BazhanovMaxim/metrics/internal/server/json"
	"github.com/BazhanovMaxim/metrics/internal/server/model"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestFileStorage_Update(t *testing.T) {
	filePath := "test.json"
	file, _ := createTestFile(filePath)
	writeToFileData(filePath, []byte("{\n    \"id\": \"CounterBatchZip55\",\n    \"type\": \"counter\",\n    \"delta\": 1353784266\n  }"))
	defer file.Close()
	defer removeTestFile(filePath)
	t.Run("Check update metric", func(t *testing.T) {
		ft := NewFileStorage(filePath, 0, false)
		_, err := ft.Update(model.Metrics{ID: "CounterBatchZip55", MType: "counter", Delta: int64Pointer(100)})
		assert.NoError(t, err)

		var metrics model.StorageJSONMetrics
		readFileAndMarshalToPojo(filePath, &metrics)

		assert.True(t, len(metrics.Gauge) == 0)
		assert.True(t, len(metrics.Counter) == 1)
		assert.True(t, metrics.Counter["CounterBatchZip55"] == 100)
	})
}

func TestFileStorage_UpdateBatches(t *testing.T) {
	filePath := "test.json"
	file, _ := createTestFile(filePath)
	writeToFileData(filePath, []byte("{\n    \"id\": \"CounterBatchZip55\",\n    \"type\": \"counter\",\n    \"delta\": 1353784266\n  }"))
	defer file.Close()
	defer removeTestFile(filePath)
	t.Run("Check update batches", func(t *testing.T) {
		ft := NewFileStorage(filePath, 0, false)

		slice := make([]model.Metrics, 2)
		slice[0] = model.Metrics{ID: "CounterBatchZip55", MType: "counter", Delta: int64Pointer(100)}
		slice[1] = model.Metrics{ID: "NewId", MType: "counter", Delta: int64Pointer(1000)}

		err := ft.UpdateBatches(slice)
		assert.NoError(t, err)

		var metrics model.StorageJSONMetrics
		readFileAndMarshalToPojo(filePath, &metrics)

		assert.True(t, len(metrics.Gauge) == 0)
		assert.True(t, len(metrics.Counter) == 2)
		assert.True(t, metrics.Counter["CounterBatchZip55"] == 100)
		assert.True(t, metrics.Counter["NewId"] == 1000)
	})
}

// Вспомогательная функция для создания указателя на int64
func int64Pointer(i int) *int64 {
	value := int64(i)
	return &value
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

func readFileAndMarshalToPojo(filePath string, res any) {
	bt, _ := os.ReadFile(filePath)
	_ = json.UnmarshalJSON(bt, res)
}
