package handlers

import (
	"bytes"
	"github.com/BazhanovMaxim/metrics/internal/server/configs"
	"github.com/BazhanovMaxim/metrics/internal/server/service"
	"github.com/BazhanovMaxim/metrics/internal/server/storage"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_UpdateHandler(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name       string
		target     string
		httpMethod string
		want       want
	}{
		{
			name:       "Positive Gauge test #1",
			target:     "/update/gauge/temperature/10",
			httpMethod: http.MethodPost,
			want:       want{code: http.StatusOK},
		},
		{
			name:       "Positive Counter test #1",
			target:     "/update/counter/temperature/10",
			httpMethod: http.MethodPost,
			want:       want{code: http.StatusOK},
		},
		{
			name:       "Negative not found metric type",
			target:     "/update/unknown/temperature/10",
			httpMethod: http.MethodPost,
			want:       want{code: http.StatusBadRequest},
		},
		{
			name:       "Negative bad request",
			target:     "/update/unknown/temperature",
			httpMethod: http.MethodPost,
			want:       want{code: http.StatusNotFound},
		},
	}
	config, _ := configs.NewConfig()
	ms := service.NewMetricService(config, storage.NewMemStorage(), nil, nil)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router := gin.Default()
			router.Handle(test.httpMethod, "/update/:metricType/:metricTitle/:metricValue", NewHandler(config, *ms).updateMetric)

			request := httptest.NewRequest(test.httpMethod, test.target, nil)
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)

			assert.Equal(t, test.want.code, recorder.Code)
		})
	}
}

func TestHandler_UpdateJsonHandler(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name       string
		target     string
		httpMethod string
		body       []byte
		want       want
	}{
		{
			name:       "Positive Gauge test #1",
			target:     "/update",
			httpMethod: http.MethodPost,
			body:       []byte("{\"id\": \"temperature\", \"type\": \"gauge\", \"value\": 10}"),
			want:       want{code: http.StatusOK},
		},
		{
			name:       "Positive Counter test #1",
			target:     "/update",
			httpMethod: http.MethodPost,
			body:       []byte("{\"id\": \"temperature\", \"type\": \"counter\", \"delta\": 10}"),
			want:       want{code: http.StatusOK},
		},
		{
			name:       "Negative not found metric type",
			target:     "/update",
			httpMethod: http.MethodPost,
			body:       []byte("{\"id\": \"temperature\", \"type\": \"unknown\", \"delta\": 10}"),
			want:       want{code: http.StatusBadRequest},
		},
		{
			name:       "Negative bad request",
			target:     "/update",
			httpMethod: http.MethodPost,
			body:       []byte("{\"id\": \"temperature\", \"type\": \"counter\"}"),
			want:       want{code: http.StatusBadRequest},
		},
	}
	config, _ := configs.NewConfig()
	ms := service.NewMetricService(config, storage.NewMemStorage(), nil, nil)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router := gin.Default()
			router.Handle(test.httpMethod, "/update", NewHandler(config, *ms).updateMetricFromJSON)

			request := httptest.NewRequest(test.httpMethod, test.target, bytes.NewReader(test.body))
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)

			assert.Equal(t, test.want.code, recorder.Code)
		})
	}
}
