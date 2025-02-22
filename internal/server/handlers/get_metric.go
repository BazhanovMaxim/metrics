package handlers

import (
	"fmt"
	"github.com/BazhanovMaxim/metrics/internal/server/model"
	"github.com/BazhanovMaxim/metrics/internal/server/service"
	"github.com/BazhanovMaxim/metrics/internal/server/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) getMetric(context *gin.Context) {
	metricType := context.Param("metricType")
	metricTitle := context.Param("metricTitle")
	if metric, ok := service.NewMetricService().GetMetricValue(metricType, metricTitle, h.storage); ok {
		if metricType == string(model.Gauge) {
			context.String(http.StatusOK, fmt.Sprintf("%g", *metric.Value))
			return
		}
		context.String(http.StatusOK, fmt.Sprintf("%d", *metric.Delta))
		return
	}
	context.Status(http.StatusNotFound)
}

func (h *Handler) getMetricFromJSON(context *gin.Context) {
	modelMetrics := model.Metrics{}
	if err := utils.MarshalRequest(context, &modelMetrics); err != nil {
		context.Status(http.StatusBadRequest)
		return
	}
	context.Header("Content-Type", "application/json")
	if value, ok := service.NewMetricService().GetMetricValue(modelMetrics.MType, modelMetrics.ID, h.storage); ok {
		context.JSON(http.StatusOK, value)
		return
	}
	context.Status(http.StatusNotFound)
}
