package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/DanilWaliev/wishlist-api/internal/models"
)

type WishlistItemsRepo interface {
	Create(ctx context.Context, item *models.WishlistItem) (uint32, error)
	ReadByID(ctx context.Context, id uint32) (*models.WishlistItem, error)
	ReadByWishlistID(ctx context.Context, wishlistID uint32) ([]*models.WishlistItem, error)
	Update(ctx context.Context, item *models.WishlistItem) error
	Delete(ctx context.Context, id uint32) error
	Reserve(ctx context.Context, id uint32) error
}

type WishlistOwnerRepo interface {
	ReadByID(ctx context.Context, id uint32) (*models.Wishlist, error)
}

type WishlistItemsService struct {
	itemsRepo     WishlistItemsRepo
	wishlistsRepo WishlistOwnerRepo
}

func NewWishlistItemsService(itemsRepo WishlistItemsRepo, wishlistsRepo WishlistOwnerRepo) *WishlistItemsService {
	return &WishlistItemsService{
		itemsRepo:     itemsRepo,
		wishlistsRepo: wishlistsRepo,
	}
}

func (s *WishlistItemsService) Create(
	ctx context.Context,
	userID uint32,
	wishlistID uint32,
	req *models.CreateWishlistItemRequest,
) (*models.WishlistItemResponse, error) {
	wishlist, err := s.wishlistsRepo.ReadByID(ctx, wishlistID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("wishlist not found")
		}
		return nil, fmt.Errorf("read wishlist: %w", err)
	}

	if wishlist.UserID != userID {
		return nil, errors.New("access denied")
	}

	item := &models.WishlistItem{
		Title:       req.Title,
		Description: req.Description,
		ProductURL:  req.ProductURL,
		Priority:    req.Priority,
		Reserved:    false,
		WishlistID:  wishlistID,
	}

	id, err := s.itemsRepo.Create(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("create wishlist item: %w", err)
	}

	item.ID = id

	return mapWishlistItemToResponse(item), nil
}

func (s *WishlistItemsService) GetByID(
	ctx context.Context,
	userID uint32,
	itemID uint32,
) (*models.WishlistItemResponse, error) {
	item, err := s.itemsRepo.ReadByID(ctx, itemID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("wishlist item not found")
		}
		return nil, fmt.Errorf("read wishlist item: %w", err)
	}

	wishlist, err := s.wishlistsRepo.ReadByID(ctx, item.WishlistID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("wishlist not found")
		}
		return nil, fmt.Errorf("read wishlist: %w", err)
	}

	if wishlist.UserID != userID {
		return nil, errors.New("access denied")
	}

	return mapWishlistItemToResponse(item), nil
}

func (s *WishlistItemsService) GetByWishlistID(
	ctx context.Context,
	userID uint32,
	wishlistID uint32,
) ([]models.WishlistItemResponse, error) {
	wishlist, err := s.wishlistsRepo.ReadByID(ctx, wishlistID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("wishlist not found")
		}
		return nil, fmt.Errorf("read wishlist: %w", err)
	}

	if wishlist.UserID != userID {
		return nil, errors.New("access denied")
	}

	items, err := s.itemsRepo.ReadByWishlistID(ctx, wishlistID)
	if err != nil {
		return nil, fmt.Errorf("read wishlist items: %w", err)
	}

	resp := make([]models.WishlistItemResponse, 0, len(items))
	for i := range items {
		resp = append(resp, *mapWishlistItemToResponse(items[i]))
	}

	return resp, nil
}

func (s *WishlistItemsService) Update(
	ctx context.Context,
	userID uint32,
	itemID uint32,
	req *models.UpdateWishlistItemRequest,
) (*models.WishlistItemResponse, error) {
	item, err := s.itemsRepo.ReadByID(ctx, itemID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("wishlist item not found")
		}
		return nil, fmt.Errorf("read wishlist item: %w", err)
	}

	wishlist, err := s.wishlistsRepo.ReadByID(ctx, item.WishlistID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("wishlist not found")
		}
		return nil, fmt.Errorf("read wishlist: %w", err)
	}

	if wishlist.UserID != userID {
		return nil, errors.New("access denied")
	}

	item.Title = req.Title
	item.Description = req.Description
	item.ProductURL = req.ProductURL
	item.Priority = req.Priority

	if err := s.itemsRepo.Update(ctx, item); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("wishlist item not found")
		}
		return nil, fmt.Errorf("update wishlist item: %w", err)
	}

	return mapWishlistItemToResponse(item), nil
}

func (s *WishlistItemsService) Delete(
	ctx context.Context,
	userID uint32,
	itemID uint32,
) error {
	item, err := s.itemsRepo.ReadByID(ctx, itemID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("wishlist item not found")
		}
		return fmt.Errorf("read wishlist item: %w", err)
	}

	wishlist, err := s.wishlistsRepo.ReadByID(ctx, item.WishlistID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("wishlist not found")
		}
		return fmt.Errorf("read wishlist: %w", err)
	}

	if wishlist.UserID != userID {
		return errors.New("access denied")
	}

	if err := s.itemsRepo.Delete(ctx, itemID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("wishlist item not found")
		}
		return fmt.Errorf("delete wishlist item: %w", err)
	}

	return nil
}

func (s *WishlistItemsService) Reserve(
	ctx context.Context,
	wishlistToken string,
	itemID uint32,
) error {
	item, err := s.itemsRepo.ReadByID(ctx, itemID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("wishlist item not found")
		}
		return fmt.Errorf("read wishlist item: %w", err)
	}

	wishlist, err := s.wishlistsRepo.ReadByID(ctx, item.WishlistID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("wishlist not found")
		}
		return fmt.Errorf("read wishlist: %w", err)
	}

	if wishlist.Token != wishlistToken {
		return errors.New("invalid wishlist token")
	}

	if err := s.itemsRepo.Reserve(ctx, itemID); err != nil {
		return fmt.Errorf("reserve wishlist item: %w", err)
	}

	return nil
}

func mapWishlistItemToResponse(item *models.WishlistItem) *models.WishlistItemResponse {
	return &models.WishlistItemResponse{
		ID:          item.ID,
		Title:       item.Title,
		Description: item.Description,
		ProductURL:  item.ProductURL,
		Priority:    item.Priority,
		Reserved:    item.Reserved,
	}
}
