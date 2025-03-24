package handlers

import (
	"github.com/BazhanovMaxim/metrics/internal/server/configs"
	"github.com/BazhanovMaxim/metrics/internal/server/middleware"
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

	router.RedirectTrailingSlash = false

	// Загрузка шаблонов
	router.LoadHTMLGlob("internal/server/templates/*")

	// Регистрация middleware
	router.Use(
		middleware.GzipCompress(),
		middleware.GzipDecompress(),
		middleware.ComputeSHA256(h.config),
		middleware.NewServerLogger(),
	)

	router.GET("/", h.homePage)
	router.GET("/ping", h.ping)
	router.GET("/value/:metricType/:metricTitle", h.getMetric)
	router.POST("/value/", h.getMetricFromJSON)
	router.POST("/update/:metricType/:metricTitle/:metricValue", h.updateMetric)
	router.POST("/update/", h.updateMetricFromJSON)
	router.POST("/updates/", h.updates)

	return http.ListenAndServe(h.config.RunAddress, router)
}
