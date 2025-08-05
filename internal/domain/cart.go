package domain

import "time"

type CartOrder struct {
	UserID    uint        `json:"user_id"`
	Total     float64     `json:"total"`
	Status    bool        `json:"status"`
	Items     []CartItems `json:"items"`
	CreatedAt time.Time   `json:"created_at"`
}

type CartItems struct {
	Quantity  int     `json:"quantity"`
	ProductId uint    `json:"product_id"`
	Fee       float64 `json:"fee"`
}
