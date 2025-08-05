package domain

import "time"

type CartItems struct {
	Quantity  int       `json:"quantity" db:"quantity"`
	ProductId uint      `json:"product_id" db:"product_id"`
	Fee       float64   `json:"fee" db:"fee"`
	UserId    uint      `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type RequestCart struct {
	UserID uint        `json:"user_id"`
	Items  []CartItems `json:"items"`
}
