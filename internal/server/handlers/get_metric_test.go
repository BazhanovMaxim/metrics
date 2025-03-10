package handlers

import (
	"bytes"
	"github.com/BazhanovMaxim/metrics/internal/server/configs"
	"github.com/BazhanovMaxim/metrics/internal/server/model"
	"github.com/BazhanovMaxim/metrics/internal/server/service"
	"github.com/BazhanovMaxim/metrics/internal/server/storage"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_GetMetric(t *testing.T) {
	type want struct {
		code     int
		response string
	}
	tests := []struct {
		name         string
		method       string
		relativePath string
		targetPath   string
		want         want
	}{
		{"Positive get gauge",
			"GET",
			"/value/:metricType/:metricTitle",
			"/value/gauge/test",
			want{code: http.StatusOK, response: "10"},
		},
		{"Positive get counter",
			"GET",
			"/value/:metricType/:metricTitle",
			"/value/gauge/test",
			want{code: http.StatusOK, response: "10"},
		},
	}
	memStorage := storage.NewMemStorage()
	_, _ = memStorage.Update(model.Metrics{ID: "test", MType: "gauge", Value: float64Pointer(10)})
	_, _ = memStorage.Update(model.Metrics{ID: "test", MType: "counter", Delta: int64Pointer(10)})

	config, _ := configs.NewConfig()
	ms := service.NewMetricService(config, memStorage)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router := gin.Default()
			router.Handle(test.method, test.relativePath, NewHandler(config, *ms).getMetric)

			request := httptest.NewRequest(test.method, test.targetPath, nil)
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)

			assert.Equal(t, test.want.code, recorder.Code)
			if recorder.Code == test.want.code {
				assert.Equal(t, test.want.response, recorder.Body.String())
			}
		})
	}
}

func TestHandler_GetMetricFromJson(t *testing.T) {
	type want struct {
		code     int
		response string
	}
	tests := []struct {
		name         string
		method       string
		relativePath string
		targetPath   string
		body         []byte
		want         want
	}{
		{"Positive get gauge",
			"GET",
			"/value",
			"/value",
			[]byte("{\"id\": \"test\", \"type\": \"gauge\"}"),
			want{code: http.StatusOK, response: "{\"id\":\"test\",\"type\":\"gauge\",\"value\":10}"},
		},
		{"Positive get counter",
			"GET",
			"/value",
			"/value",
			[]byte("{\"id\": \"test\", \"type\": \"counter\"}"),
			want{code: http.StatusOK, response: "{\"id\":\"test\",\"type\":\"counter\",\"delta\":10}"},
		},
	}
	ms := storage.NewMemStorage()
	_, _ = ms.Update(model.Metrics{ID: "test", MType: "gauge", Value: float64Pointer(10)})
	_, _ = ms.Update(model.Metrics{ID: "test", MType: "counter", Delta: int64Pointer(10)})
	config, _ := configs.NewConfig()
	serv := service.NewMetricService(config, ms)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router := gin.Default()
			router.Handle(test.method, test.relativePath, NewHandler(config, *serv).getMetricFromJSON)

			request := httptest.NewRequest(test.method, test.targetPath, bytes.NewReader(test.body))
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)

			assert.Equal(t, test.want.code, recorder.Code)
			if recorder.Code == test.want.code {
				assert.Equal(t, test.want.response, recorder.Body.String())
			}
		})
	}
}

func int64Pointer(i int) *int64 {
	value := int64(i)
	return &value
}

func float64Pointer(i float64) *float64 {
	return &i
}
