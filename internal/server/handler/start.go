package handler

import (
	"github.com/BazhanovMaxim/metrics/internal/server/repository"
	"net/http"
)

type Handler struct {
	storage repository.MetricStorage
}

func NewHandler(storage *repository.MetricStorage) *Handler {
	return &Handler{storage: *storage}
}

func (h *Handler) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", h.HomePageHandler)
	mux.HandleFunc("/update/", h.UpdateHandler)

	return http.ListenAndServe(":8080", mux)
}
