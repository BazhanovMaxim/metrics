package utils

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

func MarshalRequest(context *gin.Context, model any) error {
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(context.Request.Body); err != nil {
		return err
	}
	if err := json.Unmarshal(buf.Bytes(), model); err != nil {
		return err
	}
	return nil
}
