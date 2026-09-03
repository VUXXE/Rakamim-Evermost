package model

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	ID      uint   `json:"id"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Name    string `json:"name"`
	IsAdmin bool   `json:"is_admin"`
}

type StoreResponse struct {
	ID      uint   `json:"id"`
	UserID  uint   `json:"user_id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

type AuthResponse struct {
	Token string         `json:"token"`
	User  UserResponse   `json:"user"`
	Store *StoreResponse `json:"store,omitempty"`
}
