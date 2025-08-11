package domain

import "time"

type CartItems struct {
	Quantity  int       `json:"quantity" db:"quantity"`
	ProductId uint      `json:"product_id" db:"product_id"`
	Fee       float64   `json:"-" db:"fee"`
	UserId    uint      `json:"-" db:"user_id"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
	Discount  float64   `json:"-" db:"discount"`
}

type RequestCart struct {
	Items []CartItems `json:"items"`
}
