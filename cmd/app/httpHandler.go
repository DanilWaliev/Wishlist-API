package main

import "github.com/DanilWaliev/wishlist-api/internal/handlers"

// контейнер для всех HTTP обработчиков
type HTTPHandler struct {
	authHandler *handlers.AuthHandler
}

func NewHTTPHandler(authHandler *handlers.AuthHandler) *HTTPHandler {
	return &HTTPHandler{
		authHandler: authHandler,
	}
}
