package handlers

import (
	"fmt"
	"github.com/BazhanovMaxim/metrics/internal/agent/configs"
	"github.com/BazhanovMaxim/metrics/internal/agent/logger"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type Handler struct {
	config configs.Config
	client *resty.Client
}

func NewHandler(config configs.Config, client *resty.Client) *Handler {
	return &Handler{config: config, client: client}
}

func (h *Handler) SendMetrics(body []byte) {
	url := fmt.Sprintf("http://%s/update/", h.config.RunAddress)
	_, err := h.client.R().
		SetBody(body).
		SetHeader("Content-Type", "application/json").
		Post(url)
	if err != nil {
		logger.Log.Error("Failed to send metrics",
			zap.ByteString("body: ", body),
			zap.String("url: ", url),
			zap.Error(err),
		)
	}
}
