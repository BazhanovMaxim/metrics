package handler

import "net/http"

func (h *Handler) HomePageHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusBadRequest)
	//w.Write([]byte("Стартовая страница!"))
}
