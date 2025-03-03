package handlers

import (
	"github.com/BazhanovMaxim/metrics/internal/server/compress"
	"github.com/BazhanovMaxim/metrics/internal/server/configs"
	"github.com/BazhanovMaxim/metrics/internal/server/logger"
	"github.com/BazhanovMaxim/metrics/internal/server/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	config  configs.Config
	service service.MetricService
}

func NewHandler(config configs.Config, service service.MetricService) *Handler {
	return &Handler{config: config, service: service}
}

func (h *Handler) Start() error {
	router := gin.Default()

	// Загрузка шаблонов
	router.LoadHTMLGlob("internal/server/templates/*")

	// Регистрация middleware
	router.Use(
		compress.GzipCompress(),
		compress.GzipDecompress(),
		logger.RequestLogger(),
		logger.ResponseLogger(),
	)

	router.GET("/", h.homePage)
	router.GET("/value/:metricType/:metricTitle", h.getMetric)
	router.POST("/value", h.getMetricFromJSON)
	router.POST("/update/:metricType/:metricTitle/:metricValue", h.updateMetric)
	router.POST("/update", h.updateMetricFromJSON)

	return http.ListenAndServe(h.config.RunAddress, router)
}
