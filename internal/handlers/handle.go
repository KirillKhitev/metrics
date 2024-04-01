package handlers

import "net/http"

// BadRequestHandle возвращает код ответа 400.
func BadRequestHandle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}

// NotFoundHandle возвращает код ответа 404.
func NotFoundHandle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}
