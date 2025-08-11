package domain

import "time"

type Payment struct {
	Id        uint      `json:"-" db:"id"`
	UserID    uint      `json:"-" db:"user_id"`
	Status    string    `json:"status" db:"status"`
	OrderID   uint      `json:"-" db:"status"`
	Total     float64   `json:"total" db:"total"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}
