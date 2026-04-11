package services

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/DanilWaliev/wishlist-api/internal/models"
)

type WishlistsService struct {
	wlRepo    WishlistsRepo
	itemsRepo WishlistItemsRepo
}

func NewWishlistsService(wlRepo WishlistsRepo, itemsRepo WishlistItemsRepo) *WishlistsService {
	return &WishlistsService{
		wlRepo:    wlRepo,
		itemsRepo: itemsRepo,
	}
}

type WishlistsRepo interface {
	Create(ctx context.Context, w *models.Wishlist) (uint32, error)
	ReadByID(ctx context.Context, id uint32) (*models.Wishlist, error)
	ReadByUserID(ctx context.Context, userId uint32) ([]*models.Wishlist, error)
	Update(ctx context.Context, w *models.Wishlist) error
	Delete(ctx context.Context, id uint32) error
	ReadByToken(ctx context.Context, token string) (*models.Wishlist, error)
}

func (s *WishlistsService) Create(
	ctx context.Context,
	userID uint32,
	req *models.CreateWishlistRequest,
) (*models.WishlistResponse, error) {

	eventDate, err := time.Parse("2006-01-02", req.EventDate)
	if err != nil {
		return nil, errors.New("invalid event_date format, expected YYYY-MM-DD")
	}

	token, err := generateToken()
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	w := &models.Wishlist{
		EventName:   req.EventName,
		Description: req.Description,
		EventDate:   eventDate,
		Token:       token,
		UserID:      userID,
	}

	id, err := s.wlRepo.Create(ctx, w)
	if err != nil {
		return nil, fmt.Errorf("create wishlist: %w", err)
	}

	w.ID = id

	return mapWishlistToResponse(w), nil
}

func (s *WishlistsService) GetByID(
	ctx context.Context,
	userID uint32,
	wishlistID uint32,
) (*models.WishlistResponse, error) {

	w, err := s.wlRepo.ReadByID(ctx, wishlistID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("wishlist not found")
		}
		return nil, fmt.Errorf("read wishlist: %w", err)
	}

	if w.UserID != userID {
		return nil, errors.New("access denied")
	}

	return mapWishlistToResponse(w), nil
}

func (s *WishlistsService) GetByUserID(
	ctx context.Context,
	userID uint32,
) ([]*models.WishlistResponse, error) {

	wishlists, err := s.wlRepo.ReadByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("read wishlists: %w", err)
	}

	resp := make([]*models.WishlistResponse, 0, len(wishlists))

	for i := range wishlists {
		resp = append(resp, mapWishlistToResponse(wishlists[i]))
	}

	return resp, nil
}

func (s *WishlistsService) Update(
	ctx context.Context,
	userID uint32,
	wishlistID uint32,
	req *models.UpdateWishlistRequest,
) (*models.WishlistResponse, error) {

	w, err := s.wlRepo.ReadByID(ctx, wishlistID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("wishlist not found")
		}
		return nil, fmt.Errorf("read wishlist: %w", err)
	}

	if w.UserID != userID {
		return nil, errors.New("access denied")
	}

	eventDate, err := time.Parse("2006-01-02", req.EventDate)
	if err != nil {
		return nil, errors.New("invalid event_date format")
	}

	w.EventName = req.EventName
	w.Description = req.Description
	w.EventDate = eventDate

	if err := s.wlRepo.Update(ctx, w); err != nil {
		return nil, fmt.Errorf("update wishlist: %w", err)
	}

	return mapWishlistToResponse(w), nil
}

func (s *WishlistsService) Delete(
	ctx context.Context,
	userID uint32,
	wishlistID uint32,
) error {

	w, err := s.wlRepo.ReadByID(ctx, wishlistID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("wishlist not found")
		}
		return fmt.Errorf("read wishlist: %w", err)
	}

	if w.UserID != userID {
		return errors.New("access denied")
	}

	if err := s.wlRepo.Delete(ctx, wishlistID); err != nil {
		return fmt.Errorf("delete wishlist: %w", err)
	}

	return nil
}

func (s *WishlistsService) GetPublicByToken(
	ctx context.Context,
	token string,
) (*models.PublicWishlistResponse, error) {

	w, err := s.wlRepo.ReadByToken(ctx, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("wishlist not found")
		}
		return nil, fmt.Errorf("read wishlist by token: %w", err)
	}

	items, err := s.itemsRepo.ReadByWishlistID(ctx, w.ID)
	if err != nil {
		return nil, fmt.Errorf("read wishlist items: %w", err)
	}

	respItems := make([]models.WishlistItemResponse, 0, len(items))

	for i := range items {
		respItems = append(respItems, models.WishlistItemResponse{
			ID:          items[i].ID,
			Title:       items[i].Title,
			Description: items[i].Description,
			ProductURL:  items[i].ProductURL,
			Priority:    items[i].Priority,
			Reserved:    items[i].Reserved,
		})
	}

	return &models.PublicWishlistResponse{
		ID:          w.ID,
		EventName:   w.EventName,
		Description: w.Description,
		EventDate:   w.EventDate.Format("2006-01-02"),
		Items:       respItems,
	}, nil
}

// вспомогательные функции
func generateToken() (string, error) {
	b := make([]byte, 16)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func mapWishlistToResponse(w *models.Wishlist) *models.WishlistResponse {
	return &models.WishlistResponse{
		ID:          w.ID,
		EventName:   w.EventName,
		Description: w.Description,
		EventDate:   w.EventDate.Format("2006-01-02"),
		Token:       w.Token,
	}
}
