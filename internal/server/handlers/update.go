package handlers

import (
	"github.com/BazhanovMaxim/metrics/internal/server/model"
	"github.com/BazhanovMaxim/metrics/internal/server/service"
	"github.com/BazhanovMaxim/metrics/internal/server/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) updateMetric(context *gin.Context) {
	if context.Request.Method != http.MethodPost {
		context.Status(http.StatusMethodNotAllowed)
		return
	}
	metricType := context.Param("metricType")
	metricTitle := context.Param("metricTitle")
	metricValue := context.Param("metricValue")
	_, code := h.update(metricTitle, metricType, metricValue)
	context.Status(code)
}

func (h *Handler) updateMetricFromJSON(context *gin.Context) {
	if context.Request.Method != http.MethodPost {
		context.Status(http.StatusMethodNotAllowed)
		return
	}
	modelMetrics := model.Metrics{}
	if err := utils.MarshalRequest(context, &modelMetrics); err != nil {
		context.Status(http.StatusBadRequest)
		return
	}
	if modelMetrics.Delta == nil && modelMetrics.Value == nil {
		context.Status(http.StatusNotFound)
		return
	}
	context.Header("Content-Type", "application/json")
	if modelMetrics.MType == string(model.Counter) {
		var delta int64
		if modelMetrics.Delta == nil {
			delta = 0
		} else {
			delta = *modelMetrics.Delta
		}
		response, code := h.update(modelMetrics.ID, modelMetrics.MType, delta)
		context.JSON(code, response)
		return
	}
	var value float64
	if modelMetrics.Value == nil {
		value = 0.0
	} else {
		value = *modelMetrics.Value
	}
	response, code := h.update(modelMetrics.ID, modelMetrics.MType, value)
	context.JSON(code, response)
}

func (h *Handler) update(id, metricType string, value interface{}) (*model.Metrics, int) {
	metric := service.NewMetricService().FindService(metricType)
	if metric == nil {
		return nil, http.StatusBadRequest
	}
	updatedMetric := metric(id, value, h.storage)
	if updatedMetric == nil {
		return nil, http.StatusBadRequest
	}
	return updatedMetric, http.StatusOK
}
