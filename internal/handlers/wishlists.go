package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/DanilWaliev/wishlist-api/internal/middleware"
	"github.com/DanilWaliev/wishlist-api/internal/models"
	"github.com/DanilWaliev/wishlist-api/internal/services"
)

type WishlistsHandler struct {
	service *services.WishlistsService
}

func NewWishlistsHandler(service *services.WishlistsService) *WishlistsHandler {
	return &WishlistsHandler{
		service: service,
	}
}

func (h *WishlistsHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req models.CreateWishlistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.EventName = strings.TrimSpace(req.EventName)
	req.Description = strings.TrimSpace(req.Description)
	req.EventDate = strings.TrimSpace(req.EventDate)

	if req.EventName == "" || req.EventDate == "" {
		writeError(w, http.StatusBadRequest, "event_name and event_date are required")
		return
	}

	resp, err := h.service.Create(r.Context(), userID, &req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *WishlistsHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := parseUint32PathValue(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid wishlist id")
		return
	}

	resp, err := h.service.GetByID(r.Context(), userID, id)
	if err != nil {
		status := mapWishlistServiceError(err)
		writeError(w, status, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *WishlistsHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	resp, err := h.service.GetByUserID(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *WishlistsHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := parseUint32PathValue(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid wishlist id")
		return
	}

	var req models.UpdateWishlistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.EventName = strings.TrimSpace(req.EventName)
	req.Description = strings.TrimSpace(req.Description)
	req.EventDate = strings.TrimSpace(req.EventDate)

	if req.EventName == "" || req.EventDate == "" {
		writeError(w, http.StatusBadRequest, "event_name and event_date are required")
		return
	}

	resp, err := h.service.Update(r.Context(), userID, id, &req)
	if err != nil {
		status := mapWishlistServiceError(err)
		writeError(w, status, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *WishlistsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := parseUint32PathValue(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid wishlist id")
		return
	}

	err = h.service.Delete(r.Context(), userID, id)
	if err != nil {
		status := mapWishlistServiceError(err)
		writeError(w, status, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "wishlist deleted successfully",
	})
}

func mapWishlistServiceError(err error) int {
	switch err.Error() {
	case "wishlist not found":
		return http.StatusNotFound
	case "access denied":
		return http.StatusForbidden
	case "invalid event_date format, expected YYYY-MM-DD":
		return http.StatusBadRequest
	case "invalid event_date format":
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
