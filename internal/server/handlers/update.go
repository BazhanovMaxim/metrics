package handlers

import (
	"github.com/BazhanovMaxim/metrics/internal/server/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) UpdateHandler(context *gin.Context) {
	if context.Request.Method != "POST" {
		context.Status(http.StatusMethodNotAllowed)
		return
	}
	metricType := context.Param("metricType")
	metricTitle := context.Param("metricTitle")
	metricValue := context.Param("metricValue")
	metric, ok := service.NewMetricService().FindService(metricType)
	if !ok {
		context.Status(http.StatusBadRequest)
		return
	}
	if err := metric(metricTitle, metricValue, h.storage) != nil; err {
		context.Status(http.StatusBadRequest)
		return
	}
	context.Status(http.StatusOK)
}
