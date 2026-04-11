package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/DanilWaliev/wishlist-api/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo UserRepo
	secret   []byte
}

// интерфейс репозитория БД
type UserRepo interface {
	Create(ctx context.Context, u *models.User) error
	ReadByID(ctx context.Context, id uint32) (*models.User, error)
	ReadByEmail(ctx context.Context, email string) (*models.User, error)
}

func NewAuthService(userRepo UserRepo, secret []byte) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		secret:   secret,
	}
}

func (s *AuthService) Register(ctx context.Context, r *models.RegisterRequest) (*models.RegisterResponse, error) {
	user, err := s.userRepo.ReadByEmail(ctx, r.Email)
	if err == nil && user != nil {
		return nil, errors.New("пользователь с такой почтой уже существует")
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w", err)
	}

	hashedPwd, err := hashPassword(r.Password)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	u := &models.User{
		Email:        r.Email,
		PasswordHash: hashedPwd,
	}

	err = s.userRepo.Create(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	u, err = s.userRepo.ReadByEmail(ctx, r.Email)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &models.RegisterResponse{
		ID:    u.ID,
		Email: u.Email,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, r *models.LoginRequest) (*models.LoginResponse, error) {
	user, err := s.userRepo.ReadByEmail(ctx, r.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("неверная почта или пароль")
		}
		return nil, fmt.Errorf("%w", err)
	}

	if !checkPassword(r.Password, user.PasswordHash) {
		return nil, errors.New("неверная почта или пароль")
	}

	token, err := s.GenerateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &models.LoginResponse{
		Token: token,
	}, nil
}

func (s *AuthService) GenerateToken(userID uint32) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return signedToken, nil
}

func (s *AuthService) AuthorizeByToken(ctx context.Context, tokenString string) (*models.User, error) {
	token, err := s.parseToken(tokenString)
	if err != nil {
		return nil, errors.New("невалидный токен")
	}

	if !token.Valid {
		return nil, errors.New("невалидный токен")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("не удалось прочитать claims токена")
	}

	rawUserID, ok := claims["user_id"]
	if !ok {
		return nil, errors.New("в токене отсутствует user_id")
	}

	userIDFloat, ok := rawUserID.(float64)
	if !ok {
		return nil, errors.New("некорректный user_id в токене")
	}

	userID := uint32(userIDFloat)

	user, err := s.userRepo.ReadByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("пользователь не найден")
		}
		return nil, fmt.Errorf("%w", err)
	}

	return user, nil
}

// вспомогательные методы
func (s *AuthService) parseToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.secret), nil
	})
}

func hashPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return string(hash), nil
}

func checkPassword(pwd string, storedPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedPwd), []byte(pwd))

	if err != nil {
		return false
	}

	return true
}
