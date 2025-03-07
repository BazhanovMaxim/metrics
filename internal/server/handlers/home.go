package handlers

import (
	"github.com/BazhanovMaxim/metrics/internal/server/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) homePage(context *gin.Context) {
	var items []model.IndexHTMLModel
	for _, metric := range h.service.GetMetrics() {
		var val interface{}
		if metric.MType == string(model.Counter) {
			val = metric.Delta
		} else {
			val = metric.Value
		}
		items = append(items, model.IndexHTMLModel{Key: metric.ID, Value: val})
	}
	data := gin.H{"Metrics": items}
	context.HTML(http.StatusOK, "index.html", data)
}
