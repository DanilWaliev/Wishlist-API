package main

import "net/http"

func routes(h *HTTPHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /register", h.authHandler.Register)
	mux.HandleFunc("POST /login", h.authHandler.Login)

	mux.Handle("GET /wishlists", h.authMW.RequireAuth(http.HandlerFunc(h.wishlistsHandler.GetByUserID)))
	mux.Handle("POST /wishlists", h.authMW.RequireAuth(http.HandlerFunc(h.wishlistsHandler.Create)))
	mux.Handle("GET /wishlists/{id}", h.authMW.RequireAuth(http.HandlerFunc(h.wishlistsHandler.GetByID)))
	mux.Handle("PUT /wishlists/{id}", h.authMW.RequireAuth(http.HandlerFunc(h.wishlistsHandler.Update)))
	mux.Handle("DELETE /wishlists/{id}", h.authMW.RequireAuth(http.HandlerFunc(h.wishlistsHandler.Delete)))

	mux.Handle("GET /wishlists/{id}/items", h.authMW.RequireAuth(http.HandlerFunc(h.wlItemsHandler.GetByWishlistID)))
	mux.Handle("POST /wishlists/{id}/items", h.authMW.RequireAuth(http.HandlerFunc(h.wlItemsHandler.Create)))
	mux.Handle("GET /wishlists/{id}/items/{itemId}", h.authMW.RequireAuth(http.HandlerFunc(h.wlItemsHandler.GetByID)))
	mux.Handle("PUT /wishlists/{id}/items/{itemId}", h.authMW.RequireAuth(http.HandlerFunc(h.wlItemsHandler.Update)))
	mux.Handle("DELETE /wishlists/{id}/items/{itemId}", h.authMW.RequireAuth(http.HandlerFunc(h.wlItemsHandler.Delete)))

	mux.HandleFunc("POST /public/{token}/reserve/{itemId}", h.wlItemsHandler.Reserve)
	mux.HandleFunc("GET /public/{token}", h.wishlistsHandler.GetPublic)

	return mux
}
