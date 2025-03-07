package compress

import (
	"bufio"
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"io"
	"net"
	"net/http"
	"strings"
)

func GzipDecompress() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Проверяем, сжато ли тело запроса с помощью Gzip
		if strings.Contains(c.Request.Header.Get("Content-Encoding"), "gzip") {
			body := DecompressBody(c)
			// Заменяем тело запроса на распакованное
			c.Request.Body = io.NopCloser(strings.NewReader(string(body)))
			c.Request.Header.Del("Content-Encoding")
			c.Request.Header.Set("Content-Length", string(rune(len(body))))
		}
		c.Next()
	}
}

func DecompressBody(c *gin.Context) []byte {
	// Читаем сжатое тело запроса
	compressedBody, _ := io.ReadAll(c.Request.Body)

	// Распаковываем тело запроса
	reader, _ := gzip.NewReader(io.NopCloser(strings.NewReader(string(compressedBody))))
	defer reader.Close()

	// Читаем распакованное тело запроса
	decompressedBody, _ := io.ReadAll(reader)
	return decompressedBody
}

func GzipCompress() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Проверяем, поддерживает ли клиент сжатие Gzip
		if strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
			// Обертываем writer для сжатия Gzip
			gw := NewGzipWriter(c.Writer)
			c.Writer = gw

			// Закрываем gzip.Writer после завершения запроса
			defer gw.Close()
		}
		c.Next()
	}
}

type GzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func NewGzipWriter(w gin.ResponseWriter) *GzipWriter {
	return &GzipWriter{
		ResponseWriter: w,
		writer:         gzip.NewWriter(w),
	}
}

func (gw *GzipWriter) Write(data []byte) (int, error) {
	return gw.writer.Write(data)
}

func (gw *GzipWriter) WriteString(data string) (int, error) {
	return gw.writer.Write([]byte(data))
}

func (gw *GzipWriter) WriteHeader(statusCode int) {
	gw.Header().Del("Content-Length")
	gw.Header().Set("Content-Encoding", "gzip")
	gw.ResponseWriter.WriteHeader(statusCode)
}

func (gw *GzipWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return gw.ResponseWriter.(http.Hijacker).Hijack()
}

func (gw *GzipWriter) Close() error {
	return gw.writer.Close()
}
