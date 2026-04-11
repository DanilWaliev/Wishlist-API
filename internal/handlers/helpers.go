package handlers

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

// Отправляет Internal Server Error и логирует ошибку
func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	log.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Отправляет ошибку с указанным статусом
func ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// Обертка ClientError для отправки статуса Not Found
func NotFound(w http.ResponseWriter) {
	ClientError(w, http.StatusNotFound)
}
