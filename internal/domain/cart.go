package domain

type CartItems struct {
	UserID    uint    `json:"user_id" db:"user_id"`
	Quantity  int     `json:"quantity" db:"quantity"`
	ProductId uint    `json:"product_id" db:"product_id"`
	Fee       float64 `json:"fee" db:"fee"`
}
