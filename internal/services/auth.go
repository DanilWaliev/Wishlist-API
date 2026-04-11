package services

import "github.com/DanilWaliev/wishlist-api/internal/repository"

type AuthService struct {
	userRepo *repository.UserRepo
}

func NewAuthService(userRepo *repository.UserRepo) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}
