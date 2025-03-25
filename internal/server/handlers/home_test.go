package handlers

import (
	"github.com/BazhanovMaxim/metrics/internal/server/configs"
	"github.com/BazhanovMaxim/metrics/internal/server/service"
	"github.com/BazhanovMaxim/metrics/internal/server/storage"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_HomePageHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name         string
		method       string
		relativePath string
		targetPath   string
		want         want
	}{
		{
			name:         "Positive Get Home Page",
			method:       "GET",
			relativePath: "/",
			targetPath:   "/",
			want: want{
				code: http.StatusOK,
				response: "<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n    " +
					"<meta charset=\"UTF-8\">\n    " +
					"<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n    " +
					"<title>Metrics Page</title>\n" +
					"</head>\n" +
					"<body>\n" +
					"<ul>\n    \n" +
					"</ul>\n" +
					"</body>\n" +
					"</html>",
			},
		},
	}
	config, _ := configs.NewConfig()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router := gin.Default()
			router.LoadHTMLGlob("../templates/*")
			router.Handle(test.method, test.relativePath, NewHandler(config, *service.NewMetricService(config, storage.NewMemStorage())).homePage)

			request := httptest.NewRequest(test.method, test.targetPath, nil)
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)

			assert.Equal(t, test.want.code, recorder.Code)
			assert.Contains(t, test.want.response, recorder.Body.String())
		})
	}
}
