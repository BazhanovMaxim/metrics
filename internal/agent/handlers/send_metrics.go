package handlers

import (
	"fmt"
	"github.com/BazhanovMaxim/metrics/internal/agent/configs"
	"github.com/go-resty/resty/v2"
)

type Handler struct {
	client *resty.Client
}

func NewHandler(client *resty.Client) *Handler {
	return &Handler{client: client}
}

func (h *Handler) SendMetrics(config configs.Config, body []byte) {
	_, _ = h.client.R().
		SetBody(body).
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Content-Type", "application/json").
		Post(fmt.Sprintf("http://%s/update/", config.RunAddress))
}
