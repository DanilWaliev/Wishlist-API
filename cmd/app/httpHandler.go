package main

import (
	"github.com/DanilWaliev/wishlist-api/internal/handlers"
	"github.com/DanilWaliev/wishlist-api/internal/middleware"
)

// контейнер для всех HTTP обработчиков
type HTTPHandler struct {
	authMW           *middleware.AuthMiddleware
	authHandler      *handlers.AuthHandler
	wishlistsHandler *handlers.WishlistsHandler
	wlItemsHandler   *handlers.WishlistItemsHandler
}

func NewHTTPHandler(authMW *middleware.AuthMiddleware,
	authHandler *handlers.AuthHandler,
	wlHandler *handlers.WishlistsHandler,
	wlItemsHandler *handlers.WishlistItemsHandler) *HTTPHandler {
	return &HTTPHandler{
		authMW:           authMW,
		authHandler:      authHandler,
		wishlistsHandler: wlHandler,
		wlItemsHandler:   wlItemsHandler,
	}
}
