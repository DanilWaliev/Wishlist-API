package main

import "net/http"

func routes(h *HTTPHandler) *http.ServeMux {
	mux := http.NewServeMux()

	// тут будут эндпоинты

	return mux
}
