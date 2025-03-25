package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/BazhanovMaxim/metrics/internal/server/configs"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

// ComputeSHA256 механизм реализации подписи алгоритма SHA256
func ComputeSHA256(config configs.Config) gin.HandlerFunc {
	return func(context *gin.Context) {
		if config.SecretKey == "" {
			context.Next()
			return
		}

		if context.Request.Header.Get("HashSHA256") != "" {
			// Обработка запроса
			if !computeRequest(context, config) {
				context.AbortWithStatus(http.StatusBadRequest)
				return
			}
		}

		// Обертка для захвата тела ответа
		writer := &responseWriter{
			ResponseWriter: context.Writer,
			secretKey:      config.SecretKey,
			body:           bytes.NewBuffer(nil),
		}
		context.Writer = writer

		context.Next()

		// Проверяем, что ответ ещё не был отправлен
		if !writer.Written() {
			// Вычисляем хеш и добавляем заголовок
			hash := writer.computeHash()
			context.Writer.Header().Set("HashSHA256", hash)

			// Если тело уже записано в буфер, отправляем его
			if writer.body.Len() > 0 {
				context.Writer.Write(writer.body.Bytes())
			}
		}
	}
}

// computeRequest обрабатывает запрос
func computeRequest(context *gin.Context, config configs.Config) bool {
	contextBody := context.Request.Body
	defer contextBody.Close()

	body, _ := io.ReadAll(contextBody)

	// Восстановление тела запроса для последующего использования
	context.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	expectedHash := computeHMACSHA256(config.SecretKey, body)
	return hmac.Equal([]byte(context.Request.Header.Get("HashSHA256")), []byte(expectedHash))
}

type responseWriter struct {
	gin.ResponseWriter
	secretKey string
	body      *bytes.Buffer
	status    int
}

func (w *responseWriter) Write(data []byte) (int, error) {
	w.body.Write(data) // Захватываем тело ответа

	// Если заголовок HashSHA256 ещё не установлен, вычисляем и добавляем
	if w.Header().Get("HashSHA256") == "" {
		hash := computeHMACSHA256(w.secretKey, data)
		w.Header().Set("HashSHA256", hash)
	}

	return w.ResponseWriter.Write(data) // Отправляем оригинальный ответ
}

func (w *responseWriter) computeHash() string {
	return computeHMACSHA256(w.secretKey, w.body.Bytes())
}

func computeHMACSHA256(secretKey string, body []byte) string {
	hm := hmac.New(sha256.New, []byte(secretKey))
	hm.Write(body)
	return hex.EncodeToString(hm.Sum(nil))
}
