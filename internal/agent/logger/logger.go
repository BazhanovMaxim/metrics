package logger

import (
	"bytes"
	"compress/gzip"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"io"
	"net/http"
)

var Log = zap.NewNop()

func NewLogger(client *resty.Client) error {
	Log, _ = zap.NewProduction()
	// Настройка логирования для Resty
	client.OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
		var body []byte
		if r.Body != nil {
			// Преобразование r.Body в []byte
			switch v := r.Body.(type) {
			case []byte:
				body = v
			case string:
				body = []byte(v)
			default:
				return nil
			}

			// Проверяем заголовок Content-Encoding на наличие gzip
			if isGzipEncoded(r.Header) {
				// Развернуть тело запроса
				var err error
				body, err = ungzip(body)
				if err != nil {
					return err
				}
			}
		}

		Log.Info("HTTP Request",
			zap.String("method", r.Method),
			zap.String("url", r.URL),
			zap.Any("headers", r.Header),
			zap.String("body", string(body)),
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

// Функция для проверки заголовка Content-Encoding на наличие gzip
func isGzipEncoded(header http.Header) bool {
	return header.Get("Content-Encoding") == "gzip"
}

// Функция для развертывания данных, сжатых в gzip
func ungzip(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return io.ReadAll(reader)
}
