package file

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestReadFile(t *testing.T) {
	filePath := "test.json"
	file, _ := createTestFile(filePath)
	defer file.Close()
	defer removeTestFile(filePath)

	t.Run("Positive read file", func(t *testing.T) {
		body, err := ReadFile(filePath)
		assert.Nil(t, err)
		assert.Equal(t, 0, len(body))
	})

	t.Run("Negative read file(no file)", func(t *testing.T) {
		body, err := ReadFile("unknownFile.json")
		assert.Nil(t, body)
		assert.NotNil(t, err)
	})
}

func TestWriteFile(t *testing.T) {
	filePath := "test.json"
	file, _ := createTestFile(filePath)
	defer file.Close()
	defer removeTestFile(filePath)

	t.Run("Check write to file", func(t *testing.T) {
		err := WriteFile(file, []byte("Hello world"))
		assert.Nil(t, err)
		res, readErr := os.ReadFile(filePath)
		assert.Nil(t, readErr)
		assert.Equal(t, "Hello world", string(res))
	})
}

func createTestFile(filePath string) (*os.File, error) {
	return os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
}

func removeTestFile(filePath string) {
	_ = os.Remove(filePath)
}
