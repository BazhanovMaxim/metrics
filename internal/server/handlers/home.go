package handlers

import (
	"github.com/BazhanovMaxim/metrics/internal/server/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) homePage(context *gin.Context) {
	var items []model.IndexHTMLModel
	for key, value := range h.service.GetCounters() {
		items = append(items, model.IndexHTMLModel{Key: key, Value: value})
	}
	for key, value := range h.service.GetGauges() {
		items = append(items, model.IndexHTMLModel{Key: key, Value: value})
	}
	data := gin.H{"Metrics": items}
	context.HTML(http.StatusOK, "index.html", data)
}
