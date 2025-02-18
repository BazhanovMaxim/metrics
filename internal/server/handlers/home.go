package handlers

import (
	"github.com/BazhanovMaxim/metrics/internal/server/model"
	"github.com/BazhanovMaxim/metrics/internal/server/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) HomePageHandler(context *gin.Context) {
	var items []model.IndexHTMLModel
	for key, value := range service.NewMetricService().GetCounters(h.storage) {
		items = append(items, model.IndexHTMLModel{Key: key, Value: value})
	}
	for key, value := range service.NewMetricService().GetGauges(h.storage) {
		items = append(items, model.IndexHTMLModel{Key: key, Value: value})
	}
	data := gin.H{"Metrics": items}
	context.HTML(http.StatusOK, "index.html", data)
}
