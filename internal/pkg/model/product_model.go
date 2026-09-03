package model

import "time"

type CreateProductRequest struct {
	CategoryID  uint    `json:"category_id" form:"category_id" validate:"required"`
	Name        string  `json:"name" form:"name" validate:"required"`
	Description string  `json:"description" form:"description"`
	Price       float64 `json:"price" form:"price" validate:"required,gte=0"`
	Quantity    int     `json:"quantity" form:"quantity" validate:"required,gte=0"`
	ImageURL    string  `json:"image_url" form:"image_url"`
}

type UpdateProductRequest struct {
	CategoryID  *uint    `json:"category_id" form:"category_id"`
	Name        string   `json:"name" form:"name"`
	Description string   `json:"description" form:"description"`
	Price       *float64 `json:"price" form:"price"`
	Quantity    *int     `json:"quantity" form:"quantity"`
	ImageURL    string   `json:"image_url" form:"image_url"`
}

type ProductFilter struct {
	Search     string `query:"search"`
	CategoryID uint   `query:"category_id"`
	StoreID    uint   `query:"store_id"`
	Sort       string `query:"sort"` // price_asc, price_desc, newest
	Limit      int    `query:"limit"`
	Offset     int    `query:"offset"`
}

type ProductResponse struct {
	ID          uint              `json:"id"`
	StoreID     uint              `json:"store_id"`
	CategoryID  uint              `json:"category_id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Price       float64           `json:"price"`
	Quantity    int               `json:"quantity"`
	ImageURL    string            `json:"image_url"`
	Store       *StoreResponse    `json:"store,omitempty"`
	Category    *CategoryResponse `json:"category,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}
