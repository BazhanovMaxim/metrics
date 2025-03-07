package compress

import (
	"bytes"
	"compress/gzip"
)

func GzipCompress(data []byte) (bytes.Buffer, error) {
	// Создаем буфер для сжатия данных
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)

	// Записываем данные в gzip-поток
	if _, err := gz.Write(data); err != nil {
		return buf, err
	}
	if err := gz.Close(); err != nil {
		return buf, err
	}
	return buf, nil
}
