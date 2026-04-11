package main

import "net/http"

func routes(h *HTTPHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /register", h.authHandler.Register)
	mux.HandleFunc("POST /login", h.authHandler.Login)

	return mux
}
