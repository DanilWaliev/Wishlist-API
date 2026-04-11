package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/DanilWaliev/wishlist-api/internal/models"
)

type WishlistsRepo struct {
	db *sql.DB
}

func NewWishlistsRepo(db *sql.DB) *WishlistsRepo {
	return &WishlistsRepo{
		db: db,
	}
}

func (r *WishlistsRepo) Create(ctx context.Context, w *models.Wishlist) (uint32, error) {
	stmt := `INSERT INTO wishlists (event_name, description, event_date, token, user_id)
	VALUES($1, $2, $3, $4, $5) RETURNING id`

	var id uint32
	err := r.db.QueryRowContext(ctx, stmt, w.EventName, w.Description, w.EventDate, w.Token, w.UserID).Scan(&id)
	if err != nil {

		
		return 0, fmt.Errorf("create wishlist: %w", err)
	}

	return id, nil
}

func (r *WishlistsRepo) ReadByID(ctx context.Context, id uint32) (*models.Wishlist, error) {
	stmt := `SELECT id, event_name, description, event_date, token, user_id FROM wishlists
	WHERE id = $1`

	row := r.db.QueryRowContext(ctx, stmt, id)
	w := &models.Wishlist{}

	err := row.Scan(&w.ID, &w.EventName, &w.Description, &w.EventDate, &w.Token, &w.UserID)

	if err != nil {
		return nil, fmt.Errorf("read wishlists: %w", err)
	}

	return w, nil
}

func (r *WishlistsRepo) ReadByUserID(ctx context.Context, userID uint32) ([]*models.Wishlist, error) {
	stmt := `
		SELECT id, event_name, description, event_date, token, user_id
		FROM wishlists
		WHERE user_id = $1
		ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, stmt, userID)
	if err != nil {
		return nil, fmt.Errorf("read wishlists by user id: %w", err)
	}
	defer rows.Close()

	wishlists := make([]*models.Wishlist, 0)

	for rows.Next() {
		var w models.Wishlist

		err := rows.Scan(
			&w.ID,
			&w.EventName,
			&w.Description,
			&w.EventDate,
			&w.Token,
			&w.UserID,
		)
		if err != nil {
			return nil, fmt.Errorf("scan wishlist: %w", err)
		}

		wishlists = append(wishlists, &w)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate wishlists: %w", err)
	}

	return wishlists, nil
}

func (r *WishlistsRepo) Update(ctx context.Context, w *models.Wishlist) error {
	stmt := `
		UPDATE wishlists
		SET event_name = $1,
		    description = $2,
		    event_date = $3
		WHERE id = $4
	`

	res, err := r.db.ExecContext(
		ctx,
		stmt,
		w.EventName,
		w.Description,
		w.EventDate,
		w.ID,
	)
	if err != nil {
		return fmt.Errorf("update wishlist: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update wishlist rowsAffected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *WishlistsRepo) Delete(ctx context.Context, id uint32) error {
	stmt := `DELETE FROM wishlists WHERE id = $1`

	res, err := r.db.ExecContext(ctx, stmt, id)
	if err != nil {
		return fmt.Errorf("delete wishlist: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete wishlist rowsAffected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
