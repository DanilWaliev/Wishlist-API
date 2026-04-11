package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/DanilWaliev/wishlist-api/internal/middleware"
	"github.com/DanilWaliev/wishlist-api/internal/models"
	"github.com/DanilWaliev/wishlist-api/internal/services"
)

type WishlistItemsHandler struct {
	service *services.WishlistItemsService
}

func NewWishlistItemsHandler(service *services.WishlistItemsService) *WishlistItemsHandler {
	return &WishlistItemsHandler{
		service: service,
	}
}

func (h *WishlistItemsHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	wishlistID, err := parseUint32PathValue(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid wishlist id")
		return
	}

	var req models.CreateWishlistItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Description = strings.TrimSpace(req.Description)
	req.ProductURL = strings.TrimSpace(req.ProductURL)

	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	resp, err := h.service.Create(r.Context(), userID, wishlistID, &req)
	if err != nil {
		writeError(w, mapWishlistItemsServiceError(err), err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *WishlistItemsHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	itemID, err := parseUint32PathValue(r, "itemId")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	resp, err := h.service.GetByID(r.Context(), userID, itemID)
	if err != nil {
		writeError(w, mapWishlistItemsServiceError(err), err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *WishlistItemsHandler) GetByWishlistID(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	wishlistID, err := parseUint32PathValue(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid wishlist id")
		return
	}

	resp, err := h.service.GetByWishlistID(r.Context(), userID, wishlistID)
	if err != nil {
		writeError(w, mapWishlistItemsServiceError(err), err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *WishlistItemsHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	itemID, err := parseUint32PathValue(r, "itemId")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	var req models.UpdateWishlistItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Description = strings.TrimSpace(req.Description)
	req.ProductURL = strings.TrimSpace(req.ProductURL)

	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	resp, err := h.service.Update(r.Context(), userID, itemID, &req)
	if err != nil {
		writeError(w, mapWishlistItemsServiceError(err), err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *WishlistItemsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	itemID, err := parseUint32PathValue(r, "itemId")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	if err := h.service.Delete(r.Context(), userID, itemID); err != nil {
		writeError(w, mapWishlistItemsServiceError(err), err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "wishlist item deleted successfully",
	})
}

func (h *WishlistItemsHandler) Reserve(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimSpace(r.PathValue("token"))
	if token == "" {
		writeError(w, http.StatusBadRequest, "invalid wishlist token")
		return
	}

	itemID, err := parseUint32PathValue(r, "itemId")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	if err := h.service.Reserve(r.Context(), token, itemID); err != nil {
		writeError(w, mapWishlistItemsServiceError(err), err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "wishlist item reserved successfully",
	})
}

func mapWishlistItemsServiceError(err error) int {
	switch err.Error() {
	case "wishlist not found":
		return http.StatusNotFound
	case "wishlist item not found":
		return http.StatusNotFound
	case "access denied":
		return http.StatusForbidden
	case "invalid wishlist token":
		return http.StatusBadRequest
	default:
		return http.StatusBadRequest
	}
}
