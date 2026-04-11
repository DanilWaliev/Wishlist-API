package models

import "time"

// доменные модели приложения

type User struct {
	ID           uint32
	Email        string
	PasswordHash string
}

type Wishlist struct {
	ID          uint32
	EventName   string
	Description string
	EventDate   time.Time
	Token       string
	UserID      uint32
}

type WishlistItem struct {
	ID          uint32
	Title       string
	Description string
	ProductURL  string
	Priority    int
	Reserved    bool
	WishlistID  uint32
}
