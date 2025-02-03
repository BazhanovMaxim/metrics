package handlers

import (
	"github.com/BazhanovMaxim/metrics/internal/server/configs"
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
			want{
				code:     http.StatusOK,
				response: "10",
			},
		},
		{"Positive get counter",
			"GET",
			"/value/:metricType/:metricTitle",
			"/value/gauge/test",
			want{
				code:     http.StatusOK,
				response: "10",
			},
		},
		{"Negative unknown metric type",
			"GET",
			"/value/:metricType/:metricTitle",
			"/value/unknown/test",
			want{code: http.StatusNotFound},
		},
	}
	repository := storage.NewMetricRepository()
	repository.Gauge.Update("test", 10)
	repository.Counter.Update("test", 10)
	config, _ := configs.NewConfig()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router := gin.Default()
			router.Handle(test.method, test.relativePath, NewHandler(config, *repository).GetMetric)

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
