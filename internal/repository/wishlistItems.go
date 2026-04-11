package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/DanilWaliev/wishlist-api/internal/models"
)

type WishlistItemsRepo struct {
	db *sql.DB
}

func NewWishlistItemsRepo(db *sql.DB) *WishlistItemsRepo {
	return &WishlistItemsRepo{
		db: db,
	}
}

func (r *WishlistItemsRepo) Create(ctx context.Context, item *models.WishlistItem) (uint32, error) {
	stmt := `
		INSERT INTO wishlist_items (title, description, product_url, priority, reserved, wishlist_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var id uint32

	err := r.db.QueryRowContext(
		ctx,
		stmt,
		item.Title,
		item.Description,
		item.ProductURL,
		item.Priority,
		item.Reserved,
		item.WishlistID,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("create wishlist item: %w", err)
	}

	return id, nil
}

func (r *WishlistItemsRepo) ReadByID(ctx context.Context, id uint32) (*models.WishlistItem, error) {
	stmt := `
		SELECT id, title, description, product_url, priority, reserved, wishlist_id
		FROM wishlist_items
		WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, stmt, id)

	item := &models.WishlistItem{}
	err := row.Scan(
		&item.ID,
		&item.Title,
		&item.Description,
		&item.ProductURL,
		&item.Priority,
		&item.Reserved,
		&item.WishlistID,
	)
	if err != nil {
		return nil, fmt.Errorf("read item: %w", err)
	}

	return item, nil
}

func (r *WishlistItemsRepo) ReadByWishlistID(ctx context.Context, wishlistID uint32) ([]*models.WishlistItem, error) {
	stmt := `
		SELECT id, title, description, product_url, priority, reserved, wishlist_id
		FROM wishlist_items
		WHERE wishlist_id = $1
		ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, stmt, wishlistID)
	if err != nil {
		return nil, fmt.Errorf("read items by wishlist id: %w", err)
	}
	defer rows.Close()

	items := make([]*models.WishlistItem, 0)

	for rows.Next() {
		var item models.WishlistItem

		err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Description,
			&item.ProductURL,
			&item.Priority,
			&item.Reserved,
			&item.WishlistID,
		)
		if err != nil {
			return nil, fmt.Errorf("scan item: %w", err)
		}

		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate items: %w", err)
	}

	return items, nil
}

func (r *WishlistItemsRepo) Update(ctx context.Context, item *models.WishlistItem) error {
	stmt := `
		UPDATE wishlist_items
		SET title = $1,
		    description = $2,
		    product_url = $3,
		    priority = $4,
		    reserved = $5
		WHERE id = $6
	`

	res, err := r.db.ExecContext(
		ctx,
		stmt,
		item.Title,
		item.Description,
		item.ProductURL,
		item.Priority,
		item.Reserved,
		item.ID,
	)
	if err != nil {
		return fmt.Errorf("update item: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update item rowsAffected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *WishlistItemsRepo) Delete(ctx context.Context, id uint32) error {
	stmt := `DELETE FROM wishlist_items WHERE id = $1`

	res, err := r.db.ExecContext(ctx, stmt, id)
	if err != nil {
		return fmt.Errorf("delete item: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete item rowsAffected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *WishlistItemsRepo) Reserve(ctx context.Context, id uint32) error {
	stmt := `
		UPDATE wishlist_items
		SET reserved = true
		WHERE id = $1 AND reserved = false
	`

	res, err := r.db.ExecContext(ctx, stmt, id)
	if err != nil {
		return fmt.Errorf("reserve wishlist item: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("reserve rowsAffected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("item already reserved or not found")
	}

	return nil
}
