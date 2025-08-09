package domain

import "time"

type Payment struct {
	ID                   uint      `json:"id" db:"id"`
	UserID               uint      `json:"user_id" db:"user_id"`
	Amount               float64   `json:"amount" db:"amount"`
	Status               string    `json:"status" db:"status"` // pending, completed, failed, cancelled
	PaymentMethod        string    `json:"payment_method" db:"payment_method"`
	GatewayTransactionID string    `json:"gateway_transaction_id" db:"gateway_transaction_id"`
	GatewayResponse      string    `json:"gateway_response" db:"gateway_response"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

type PaymentRequest struct {
	PaymentMethod string `json:"payment_method" binding:"required"` // credit_card, paypal, bank_transfer, etc.
}

type PaymentGatewayResponse struct {
	Success       bool   `json:"success"`
	TransactionID string `json:"transaction_id"`
	Message       string `json:"message"`
}

// Cart response models
type CartItemResponse struct {
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	Total       float64 `json:"total"`
}

type CartResponse struct {
	Items []CartItemResponse `json:"items"`
	Total float64            `json:"total"`
}
