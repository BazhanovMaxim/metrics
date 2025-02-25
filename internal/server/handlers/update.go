package handlers

import (
	"github.com/BazhanovMaxim/metrics/internal/server/json"
	"github.com/BazhanovMaxim/metrics/internal/server/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) updateMetric(context *gin.Context) {
	metricType := context.Param("metricType")
	metricTitle := context.Param("metricTitle")
	metricValue := context.Param("metricValue")
	_, code := h.update(metricTitle, metricType, metricValue)
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
		response, code := h.update(modelMetrics.ID, modelMetrics.MType, *modelMetrics.Delta)
		context.JSON(code, response)
		return
	}
	if modelMetrics.Value == nil {
		context.Status(http.StatusBadRequest)
		return
	}
	response, code := h.update(modelMetrics.ID, modelMetrics.MType, *modelMetrics.Value)
	context.JSON(code, response)
}

func (h *Handler) update(id, metricType string, value interface{}) (*model.Metrics, int) {
	metric := h.service.FindService(metricType)
	if metric == nil {
		return nil, http.StatusBadRequest
	}
	updatedMetric := metric(id, value)
	if updatedMetric == nil {
		return nil, http.StatusBadRequest
	}
	if h.config.StoreInterval == 0 {
		h.service.SaveMetricToStorage(updatedMetric)
	}
	return updatedMetric, http.StatusOK
}
