package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) ping(context *gin.Context) {
	if err := h.service.PingConnection(); err != nil {
		context.Status(http.StatusInternalServerError)
		return
	}
	context.Status(http.StatusOK)
}
