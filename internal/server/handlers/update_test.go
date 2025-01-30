package handlers

import (
	"github.com/BazhanovMaxim/metrics/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_UpdateHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
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
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.httpMethod, test.target, nil)
			recorder := httptest.NewRecorder()
			NewHandler(storage.NewMetricRepository()).UpdateHandler(recorder, request)

			res := recorder.Result()
			// проверяем код ответа
			res.Body.Close()
			assert.Equal(t, test.want.code, res.StatusCode)
		})
	}
}
