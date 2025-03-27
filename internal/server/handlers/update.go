package handlers

import (
	"github.com/BazhanovMaxim/metrics/internal/agent/logger"
	"github.com/BazhanovMaxim/metrics/internal/server/json"
	"github.com/BazhanovMaxim/metrics/internal/server/model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func (h *Handler) updateMetric(context *gin.Context) {
	metricType := context.Param("metricType")
	metricTitle := context.Param("metricTitle")
	metricValue := context.Param("metricValue")
	metric := model.Metrics{ID: metricTitle, MType: metricType}
	switch metricType {
	case string(model.Counter):
		val, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			context.Status(http.StatusBadRequest)
			return
		}
		metric.Delta = &val
	case string(model.Gauge):
		val, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			context.Status(http.StatusBadRequest)
			return
		}
		metric.Value = &val
	default:
		context.Status(http.StatusBadRequest)
		return
	}
	code, _ := h.update(metric)
	context.Status(code)
}

func (h *Handler) updateMetricFromJSON(context *gin.Context) {
	modelMetrics := model.Metrics{}
	if err := json.MarshalJSON(context.Request.Body, &modelMetrics); err != nil {
		context.Status(http.StatusBadRequest)
		return
	}
	context.Header("Content-Type", "application/json")
	if modelMetrics.MType == string(model.Counter) {
		if modelMetrics.Delta == nil {
			context.Status(http.StatusBadRequest)
			return
		}
		context.JSON(h.update(model.Metrics{
			ID:    modelMetrics.ID,
			MType: modelMetrics.MType,
			Delta: modelMetrics.Delta,
		}))
		return
	}
	if modelMetrics.Value == nil {
		context.Status(http.StatusBadRequest)
		return
	}
	context.JSON(h.update(model.Metrics{ID: modelMetrics.ID, MType: modelMetrics.MType, Value: modelMetrics.Value}))
}

func (h *Handler) updates(context *gin.Context) {
	var modelMetrics []model.Metrics
	if err := json.MarshalJSON(context.Request.Body, &modelMetrics); err != nil {
		logger.Log.Error("Failed to marshalJSON", zap.Error(err))
		context.Status(http.StatusBadRequest)
		return
	}
	if err := h.service.UpdateBatches(modelMetrics); err != nil {
		logger.Log.Error("Failed to update batches", zap.Error(err))
		context.Status(http.StatusBadRequest)
		return
	}
	context.Status(http.StatusOK)
}

func (h *Handler) update(m model.Metrics) (int, *model.Metrics) {
	updatedMetric, err := h.service.UpdateMetric(m)
	if err != nil {
		return http.StatusBadRequest, nil
	}
	return http.StatusOK, updatedMetric
}
