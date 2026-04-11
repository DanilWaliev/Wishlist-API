package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/DanilWaliev/wishlist-api/internal/models"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) Create(ctx context.Context, u *models.User) error {
	stmt := `INSERT INTO users (email, password_hash)
	 VALUES($1, $2)`

	_, err := r.db.Exec(stmt, u.Email, u.PasswordHash)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (r *UserRepo) ReadByID(ctx context.Context, id uint32) (*models.User, error) {
	stmt := `SELECT * FROM users
	WHERE id = $1`

	row := r.db.QueryRow(stmt, id)
	u := &models.User{}

	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash)

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return u, nil
}

func (r *UserRepo) ReadByEmail(ctx context.Context, email string) (*models.User, error) {
	stmt := `SELECT * FROM users
	WHERE email = $1`

	row := r.db.QueryRow(stmt, email)
	u := &models.User{}

	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash)

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return u, nil
}
