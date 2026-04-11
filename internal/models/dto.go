package models

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type CreateWishlistRequest struct {
	EventName   string `json:"event_name"`
	Description string `json:"description"`
	EventDate   string `json:"event_date"`
}

type UpdateWishlistRequest struct {
	EventName   string `json:"event_name"`
	Description string `json:"description"`
	EventDate   string `json:"event_date"`
}

type WishlistResponse struct {
	ID          int64                  `json:"id"`
	EventName   string                 `json:"event_name"`
	Description string                 `json:"description"`
	EventDate   string                 `json:"event_date"`
	Token       string                 `json:"token"`
	Items       []WishlistItemResponse `json:"items"`
}

type WishlistListResponse struct {
	Wishlists []WishlistResponse `json:"wishlists"`
}

type CreateWishlistItemRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ProductURL  string `json:"product_url"`
	Priority    int    `json:"priority"`
}

type UpdateWishlistItemRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ProductURL  string `json:"product_url"`
	Priority    int    `json:"priority"`
}

type WishlistItemResponse struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ProductURL  string `json:"product_url"`
	Priority    int    `json:"priority"`
	Reserved    bool   `json:"reserved"`
}

type PublicWishlistResponse struct {
	ID          int64                  `json:"id"`
	EventName   string                 `json:"event_name"`
	Description string                 `json:"description"`
	EventDate   string                 `json:"event_date"`
	Items       []WishlistItemResponse `json:"items"`
}

type ReserveItemResponse struct {
	Message string `json:"message"`
	ItemID  int64  `json:"item_id"`
}
