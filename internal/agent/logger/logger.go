package logger

import (
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

var Log = zap.NewNop()

func NewLogger(client *resty.Client) error {
	Log, _ = zap.NewProduction()
	// Настройка логирования для Resty
	client.OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
		Log.Info("HTTP Request",
			zap.String("method", r.Method),
			zap.String("url", r.URL),
			zap.Any("headers", r.Header),
			zap.Any("body", r.Body),
		)
		return nil
	})

	client.OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
		Log.Info("HTTP Response",
			zap.Int("status", r.StatusCode()),
			zap.String("status_text", r.Status()),
			zap.Any("headers", r.Header()),
			zap.String("body", string(r.Body())),
		)
		return nil
	})
	return nil
}
