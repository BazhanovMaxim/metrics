package handlers

import (
	"github.com/BazhanovMaxim/metrics/internal/server/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) GetMetric(context *gin.Context) {
	metricType := context.Param("metricType")
	metricTitle := context.Param("metricTitle")
	if value, ok := service.NewMetricService().GetMetricValue(metricType, metricTitle, h.storage); ok {
		context.String(http.StatusOK, value)
		return
	}
	context.Status(http.StatusNotFound)
}
