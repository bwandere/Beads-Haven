package models

import (
	"encoding/json"
	"time"
)

type ProductImage struct {
	ID           int    `json:"id"`
	ProductID    string `json:"product_id,omitempty"`
	URL          string `json:"url"`
	DisplayOrder int    `json:"display_order"`
}

type Product struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Brand       string         `json:"brand"`
	Type        string         `json:"type"` // bracelets, necklaces, rings, bags, fashion-accessories, art-decor
	Category    string         `json:"category"`
	Price       float64        `json:"price"`
	OldPrice    *float64       `json:"oldPrice,omitempty"`
	Rating      float64        `json:"rating"`
	Reviews     int            `json:"reviews"`
	Description string         `json:"description"`
	Tags        []string       `json:"tags"`
	Images      []string       `json:"images"`
	ImageItems  []ProductImage `json:"image_items"`
	Sizes       []string       `json:"sizes"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (p *Product) UnmarshalJSON(data []byte) error {
	type Alias Product
	aux := &struct {
		OldPriceSnake *float64 `json:"old_price"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.OldPriceSnake != nil && p.OldPrice == nil {
		p.OldPrice = aux.OldPriceSnake
	}
	return nil
}

type OrderItem struct {
	ID        int     `json:"id,omitempty"`
	OrderID   string  `json:"order_id,omitempty"`
	ProductID string  `json:"product_id"`
	Size      string  `json:"size"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type Order struct {
	ID               string      `json:"id"`
	CustomerName     string      `json:"customer_name"`
	CustomerPhone    string      `json:"customer_phone"`
	CustomerEmail    string      `json:"customer_email"`
	DeliveryLocation string      `json:"delivery_location"`
	PaymentMethod    string      `json:"payment_method"`
	MpesaPhone       string      `json:"mpesa_phone,omitempty"`
	Subtotal         float64     `json:"subtotal"`
	DeliveryFee      float64     `json:"delivery_fee"`
	TotalAmount      float64     `json:"total_amount"`
	Status           string      `json:"status"`
	CreatedAt        time.Time   `json:"created_at"`
	Items            []OrderItem `json:"items,omitempty"`
}

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"`
	Email string `json:"email"`
}
