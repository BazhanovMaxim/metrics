package handlers

import (
	"fmt"
	"github.com/BazhanovMaxim/metrics/internal/server/json"
	"github.com/BazhanovMaxim/metrics/internal/server/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) getMetric(context *gin.Context) {
	metricType := context.Param("metricType")
	metricTitle := context.Param("metricTitle")
	if metric := h.service.GetMetricValue(metricType, metricTitle); metric != nil {
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
	if err := json.MarshalJSON(context.Request.Body, &modelMetrics); err != nil {
		context.Status(http.StatusBadRequest)
		return
	}
	context.Header("Content-Type", "application/json")
	if value := h.service.GetMetricValue(modelMetrics.MType, modelMetrics.ID); value != nil {
		context.JSON(http.StatusOK, value)
		return
	}
	context.Status(http.StatusNotFound)
}
