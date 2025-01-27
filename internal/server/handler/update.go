package handler

import (
	"github.com/BazhanovMaxim/metrics/internal/server/service"
	"net/http"
	"strings"
)

func (h *Handler) UpdateHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		// разрешаем только POST-запросы
		response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	pathParts := strings.Split(request.URL.Path, "/")
	if len(pathParts) != 5 {
		response.WriteHeader(http.StatusNotFound)
		return
	}
	metricType := pathParts[2]
	metric, ok := service.GetMetricService().FindService(metricType)
	if !ok {
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := metric(pathParts, h.storage) != nil; err {
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	response.WriteHeader(http.StatusOK)
}
