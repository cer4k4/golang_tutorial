package domain

import "time"

type Payment struct {
	Id        uint      `json:"id" db:"id"`
	UserID    uint      `db:"user_id"`
	Status    string    `json:"status" db:"status"`
	OrderID   uint      `db:"status"`
	Total     float64   `json:"total" db:"total"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
