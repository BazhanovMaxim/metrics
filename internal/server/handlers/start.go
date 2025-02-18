package handlers

import (
	"github.com/BazhanovMaxim/metrics/internal/server/configs"
	"github.com/BazhanovMaxim/metrics/internal/server/logger"
	"github.com/BazhanovMaxim/metrics/internal/server/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	storage storage.MetricStorage
	config  configs.Config
}

func NewHandler(config configs.Config, storage storage.MetricStorage) *Handler {
	return &Handler{config: config, storage: storage}
}

func (h *Handler) Start() error {
	router := gin.Default()

	// Загрузка шаблонов
	router.LoadHTMLGlob("internal/server/templates/*")

	// Регистрация middleware
	router.Use(logger.RequestLoggerMiddleware(), logger.ResponseLoggerMiddleware())

	router.GET("/", h.HomePageHandler)
	router.GET("/value/:metricType/:metricTitle", h.GetMetric)
	router.POST("/update/:metricType/:metricTitle/:metricValue", h.UpdateHandler)

	return http.ListenAndServe(h.config.RunAddress, router)
}
