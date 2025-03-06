package handlers

import (
	"github.com/BazhanovMaxim/metrics/internal/server/json"
	"github.com/BazhanovMaxim/metrics/internal/server/model"
	"github.com/gin-gonic/gin"
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
	_, code := h.update(metric)
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
		response, code := h.update(model.Metrics{
			ID:    modelMetrics.ID,
			MType: modelMetrics.MType,
			Delta: modelMetrics.Delta,
		})
		context.JSON(code, response)
		return
	}
	if modelMetrics.Value == nil {
		context.Status(http.StatusBadRequest)
		return
	}
	response, code := h.update(model.Metrics{
		ID:    modelMetrics.ID,
		MType: modelMetrics.MType,
		Value: modelMetrics.Value,
	})
	context.JSON(code, response)
}

func (h *Handler) update(m model.Metrics) (*model.Metrics, int) {
	updatedMetric, err := h.service.UpdateMetric(m)
	if err != nil {
		return nil, http.StatusBadRequest
	}
	return updatedMetric, http.StatusOK
}
