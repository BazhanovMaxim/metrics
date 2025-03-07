package handlers

//
//import (
//	"bytes"
//	"github.com/BazhanovMaxim/metrics/internal/server/configs"
//	"github.com/BazhanovMaxim/metrics/internal/server/service"
//	"github.com/BazhanovMaxim/metrics/internal/server/storage"
//	"github.com/gin-gonic/gin"
//	"github.com/stretchr/testify/assert"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//)
//
//func TestHandler_GetMetric(t *testing.T) {
//	type want struct {
//		code     int
//		response string
//	}
//	tests := []struct {
//		name         string
//		method       string
//		relativePath string
//		targetPath   string
//		want         want
//	}{
//		{"Positive get gauge",
//			"GET",
//			"/value/:metricType/:metricTitle",
//			"/value/gauge/test",
//			want{code: http.StatusOK, response: "10"},
//		},
//		{"Positive get counter",
//			"GET",
//			"/value/:metricType/:metricTitle",
//			"/value/gauge/test",
//			want{code: http.StatusOK, response: "10"},
//		},
//		{"Negative unknown metric type",
//			"GET",
//			"/value/:metricType/:metricTitle",
//			"/value/unknown/test",
//			want{code: http.StatusNotFound},
//		},
//	}
//	memStorage := storage.NewMemStorage()
//	memStorage.UpdateGauge("test", 10)
//	memStorage.UpdateCounter("test", 10)
//
//	config, _ := configs.NewConfig()
//	ms := service.NewMetricService(config, memStorage, nil, nil)
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			router := gin.Default()
//			router.Handle(test.method, test.relativePath, NewHandler(config, *ms).getMetric)
//
//			request := httptest.NewRequest(test.method, test.targetPath, nil)
//			recorder := httptest.NewRecorder()
//
//			router.ServeHTTP(recorder, request)
//
//			assert.Equal(t, test.want.code, recorder.Code)
//			if recorder.Code == test.want.code {
//				assert.Equal(t, test.want.response, recorder.Body.String())
//			}
//		})
//	}
//}
//
//func TestHandler_GetMetricFromJson(t *testing.T) {
//	type want struct {
//		code     int
//		response string
//	}
//	tests := []struct {
//		name         string
//		method       string
//		relativePath string
//		targetPath   string
//		body         []byte
//		want         want
//	}{
//		{"Positive get gauge",
//			"GET",
//			"/value",
//			"/value",
//			[]byte("{\"id\": \"test\", \"type\": \"gauge\"}"),
//			want{code: http.StatusOK, response: "{\"id\":\"test\",\"type\":\"gauge\",\"value\":10}"},
//		},
//		{"Positive get counter",
//			"GET",
//			"/value",
//			"/value",
//			[]byte("{\"id\": \"test\", \"type\": \"counter\"}"),
//			want{code: http.StatusOK, response: "{\"id\":\"test\",\"type\":\"counter\",\"delta\":10}"},
//		},
//		{"Negative unknown metric type",
//			"GET",
//			"/value",
//			"/value",
//			[]byte("{\"id\": \"test\", \"type\": \"unknown\"}"),
//			want{code: http.StatusNotFound},
//		},
//	}
//	ms := storage.NewMemStorage()
//	ms.UpdateGauge("test", 10)
//	ms.UpdateCounter("test", 10)
//	config, _ := configs.NewConfig()
//	serv := service.NewMetricService(config, ms, nil, nil)
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			router := gin.Default()
//			router.Handle(test.method, test.relativePath, NewHandler(config, *serv).getMetricFromJSON)
//
//			request := httptest.NewRequest(test.method, test.targetPath, bytes.NewReader(test.body))
//			recorder := httptest.NewRecorder()
//
//			router.ServeHTTP(recorder, request)
//
//			assert.Equal(t, test.want.code, recorder.Code)
//			if recorder.Code == test.want.code {
//				assert.Equal(t, test.want.response, recorder.Body.String())
//			}
//		})
//	}
//}
