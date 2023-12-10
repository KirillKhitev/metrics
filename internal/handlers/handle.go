package handlers

import "net/http"

func BadRequestHandle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}

func NotFoundHandle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}
