package repository

import "shop/internal/domain"

type UserRepository interface {
	Create(user *domain.User) error
	GetByID(id uint) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id uint) error
}

type ProductRepository interface {
	Create(product *domain.Product) error
	GetByID(id uint) (*domain.Product, error)
	GetAll(limit, offset int) ([]*domain.Product, error)
	Update(product *domain.Product) error
	Delete(id uint) error
	UpdateStock(id uint, stock int) error
}

type OrderRepository interface {
	Create(order *domain.Order) error
	GetByID(id uint) (*domain.Order, error)
	GetByUserID(userID uint, limit, offset int) ([]*domain.Order, error)
	Update(order *domain.Order) error
	CreateOrderItem(item *domain.OrderItem) error
	GetOrderItems(orderID uint) ([]*domain.OrderItem, error)
}

type CartRepository interface {
	CreateCart(order *domain.CartOrder) error
	GetByUserID(userID uint) (*domain.CartOrder, error)
	Update(order *domain.CartOrder) error
	CreateCartItems(user uint, item *domain.CartOrder) error
	GetCartItems(orderID uint) ([]*domain.CartOrder, error)
}
