package router

import (
	"fmt"
	"github.com/BazhanovMaxim/metrics/internal/agent/configs"
	"github.com/go-resty/resty/v2"
)

type Router struct {
	client *resty.Client
}

// NewRouter создает и возвращает новый экземпляр Router
func NewRouter(client *resty.Client) *Router {
	return &Router{client: client}
}

// SendMetrics отправляет с клиента запросы на сервер
func (h *Router) SendMetrics(config configs.Config, body []byte) error {
	_, err := h.client.R().
		SetBody(body).
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Content-Type", "application/json").
		Post(fmt.Sprintf("http://%s/updates/", config.RunAddress))
	return err
}
