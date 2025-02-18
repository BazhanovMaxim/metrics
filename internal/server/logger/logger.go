package logger

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

// ResponseWriter — обертка для gin.ResponseWriter, которая позволяет захватывать тело ответа.
type ResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		start := time.Now() // точка, где выполняется хендлер
		context.Next()      // обслуживание оригинального запроса

		duration := time.Since(start) // время выполнения запроса

		// отправляем сведения о запросе в zap
		Log.Info("Http request: ",
			zap.String("method", context.Request.Method),
			zap.String("path", context.Request.URL.Path),
			zap.Duration("duration", duration),
		)
	}
}

func ResponseLoggerMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Next()

		// Создаем обертку для ResponseWriter
		writer := &ResponseWriter{body: bytes.NewBufferString(""), ResponseWriter: context.Writer}
		context.Writer = writer

		Log.Info("Http Response: ",
			zap.Int("status code", writer.Status()),
			zap.Int("body size", writer.Size()),
		)
	}
}
