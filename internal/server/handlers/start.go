package handlers

import (
	"github.com/BazhanovMaxim/metrics/internal/server/storage"
	"net/http"
)

type Handler struct {
	storage storage.IMetricStorage
}

func NewHandler(storage storage.IMetricStorage) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", h.HomePageHandler)
	mux.HandleFunc("/update/", h.UpdateHandler)

	return http.ListenAndServe(":8080", mux)
}
