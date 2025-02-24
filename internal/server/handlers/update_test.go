package handlers

import (
	"bytes"
	"github.com/BazhanovMaxim/metrics/internal/server/configs"
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
			name:       "Negative method not allowed (GET)",
			target:     "/update/unknown/temperature/10",
			httpMethod: http.MethodGet,
			want:       want{code: http.StatusMethodNotAllowed},
		},
		{
			name:       "Negative method not allowed (PUT)",
			target:     "/update/unknown/temperature/10",
			httpMethod: http.MethodPut,
			want:       want{code: http.StatusMethodNotAllowed},
		},
		{
			name:       "Negative method not allowed (DELETE)",
			target:     "/update/unknown/temperature/10",
			httpMethod: http.MethodDelete,
			want:       want{code: http.StatusMethodNotAllowed},
		},
		{
			name:       "Negative method not allowed (HEAD)",
			target:     "/update/unknown/temperature/10",
			httpMethod: http.MethodHead,
			want:       want{code: http.StatusMethodNotAllowed},
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
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router := gin.Default()
			router.Handle(test.httpMethod, "/update/:metricType/:metricTitle/:metricValue", NewHandler(config, *storage.NewMetricRepository()).updateMetric)

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
			name:       "Negative method not allowed (GET)",
			target:     "/update",
			httpMethod: http.MethodGet,
			body:       []byte("{\"id\": \"temperature\", \"type\": \"counter\", \"delta\": 10}"),
			want:       want{code: http.StatusMethodNotAllowed},
		},
		{
			name:       "Negative method not allowed (PUT)",
			target:     "/update",
			httpMethod: http.MethodPut,
			body:       []byte("{\"id\": \"temperature\", \"type\": \"counter\", \"delta\": 10}"),
			want:       want{code: http.StatusMethodNotAllowed},
		},
		{
			name:       "Negative method not allowed (DELETE)",
			target:     "/update",
			httpMethod: http.MethodDelete,
			body:       []byte("{\"id\": \"temperature\", \"type\": \"counter\", \"delta\": 10}"),
			want:       want{code: http.StatusMethodNotAllowed},
		},
		{
			name:       "Negative method not allowed (HEAD)",
			target:     "/update",
			httpMethod: http.MethodHead,
			body:       []byte("{\"id\": \"temperature\", \"type\": \"counter\", \"delta\": 10}"),
			want:       want{code: http.StatusMethodNotAllowed},
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
			want:       want{code: http.StatusNotFound},
		},
	}
	config, _ := configs.NewConfig()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router := gin.Default()
			router.Handle(test.httpMethod, "/update", NewHandler(config, *storage.NewMetricRepository()).updateMetricFromJSON)

			request := httptest.NewRequest(test.httpMethod, test.target, bytes.NewReader(test.body))
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)

			assert.Equal(t, test.want.code, recorder.Code)
		})
	}
}
