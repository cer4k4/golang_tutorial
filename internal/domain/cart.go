package domain

import "time"

type CartItems struct {
	Quantity  int       `json:"quantity" db:"quantity"`
	ProductId uint      `json:"product_id" db:"product_id"`
	Fee       float64   `db:"fee"`
	UserId    uint      `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type RequestCart struct {
	Items []CartItems `json:"items"`
}
