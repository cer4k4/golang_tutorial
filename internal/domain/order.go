package domain

import "time"

type Order struct {
    ID         uint        `json:"id" db:"id"`
    UserID     uint        `json:"user_id" db:"user_id"`
    Total      float64     `json:"total" db:"total"`
    Status     string      `json:"status" db:"status"`
    CreatedAt  time.Time   `json:"created_at" db:"created_at"`
    UpdatedAt  time.Time   `json:"updated_at" db:"updated_at"`
    Items      []OrderItem `json:"items"`
}

type OrderItem struct {
    ID        uint    `json:"id" db:"id"`
    OrderID   uint    `json:"order_id" db:"order_id"`
    ProductID uint    `json:"product_id" db:"product_id"`
    Quantity  int     `json:"quantity" db:"quantity"`
    Price     float64 `json:"price" db:"price"`
    Product   Product `json:"product"`
}

type OrderRequest struct {
    Items []OrderItemRequest `json:"items" binding:"required,dive"`
}

type OrderItemRequest struct {
    ProductID uint `json:"product_id" binding:"required"`
    Quantity  int  `json:"quantity" binding:"required,gt=0"`
}