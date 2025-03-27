package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"strings"
	"time"
)

var Log = zap.NewNop()

// NewLogger инициализирует синглтон логера с необходимым уровнем логирования.
func NewLogger(level string) error {
	// преобразуем текстовый уровень логирования в zap.AtomicLevel
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	// создаём новую конфигурацию логера
	cfg := zap.NewProductionConfig()
	// устанавливаем уровень
	cfg.Level = lvl
	// создаём логер на основе конфигурации
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	// устанавливаем синглтон
	Log = zl
	return nil
}

// ResponseWriter — обертка для gin. ResponseWriter, которая позволяет захватывать тело ответа.
type ResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func NewServerLogger() gin.HandlerFunc {
	return func(context *gin.Context) {
		// Обслуживание запроса к серверу
		start, body := requestLogger(context)

		// Обслуживание оригинального запроса
		context.Next()

		duration := time.Since(start)

		headers := make(map[string][]string)
		for name, value := range context.Request.Header {
			headers[name] = value
		}

		Log.Info("Http Request: ",
			zap.String("method", context.Request.Method),
			zap.String("path", context.Request.URL.Path),
			zap.Any("headers", headers),
			zap.ByteString("body", body),
			zap.Duration("duration", duration),
		)

		// Обслуживание ответа от сервера
		responseLogger(context)
	}
}

func requestLogger(context *gin.Context) (time.Time, []byte) {
	start := time.Now()

	// Сохраняем тело запроса
	var body []byte
	if strings.Contains(context.Request.Header.Get("Content-Encoding"), "gzip") {
		body = DecompressBody(context)
	} else {
		body, _ = io.ReadAll(context.Request.Body)
	}
	context.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	return start, body
}

func responseLogger(context *gin.Context) {
	// Создаем обертку для ResponseWriter
	writer := &ResponseWriter{body: bytes.NewBufferString(""), ResponseWriter: context.Writer}
	context.Writer = writer

	headers := make(map[string][]string)
	for name, value := range context.Writer.Header() {
		headers[name] = value
	}

	Log.Info("Http Response: ",
		zap.Int("status code", writer.Status()),
		zap.Any("headers", headers),
		zap.Int("body size", writer.Size()),
	)
}
