package handlers

import (
	"github.com/BazhanovMaxim/metrics/internal/server/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	storage storage.IMetricStorage
}

func NewHandler(storage storage.IMetricStorage) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) Start() error {
	router := gin.Default()

	// Загрузка шаблонов
	router.LoadHTMLGlob("internal/server/templates/*")

	router.GET("/", h.HomePageHandler)
	router.GET("/value/:metricType/:metricTitle", h.GetMetric)
	router.POST("/update/:metricType/:metricTitle/:metricValue", h.UpdateHandler)

	return http.ListenAndServe(":8080", router)
}
