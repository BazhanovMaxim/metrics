package router

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/BazhanovMaxim/metrics/internal/agent/compress"
	"github.com/BazhanovMaxim/metrics/internal/agent/configs"
	"github.com/BazhanovMaxim/metrics/internal/agent/logger"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
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
	request := h.client.R().
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Content-Type", "application/json")
	h.addKeyEncryption(config, body, request)
	request.SetBody(h.compressBody(body))
	_, err := request.Post(fmt.Sprintf("http://%s/updates/", config.RunAddress))
	return err
}

func (h *Router) addKeyEncryption(config configs.Config, body []byte, request *resty.Request) {
	if config.SecretKey == "" {
		return
	}
	hm := hmac.New(sha256.New, []byte(config.SecretKey))
	hm.Write(body)
	request.SetHeader("HashSHA256", hex.EncodeToString(hm.Sum(nil)))
}

func (h *Router) compressBody(body []byte) []byte {
	buf, compressErr := compress.GzipCompress(body)
	if compressErr != nil {
		logger.Log.Error("Failed to compress data", zap.Error(compressErr))
		return nil
	}
	return buf.Bytes()
}
