package usecase

import (
	"errors"
	"shop/internal/domain"
	"shop/internal/repository"
)

type OrderUsecase struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
}

func NewOrderUsecase(orderRepo repository.OrderRepository, productRepo repository.ProductRepository) *OrderUsecase {
	return &OrderUsecase{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (u *OrderUsecase) CreateOrder(userID uint, req *domain.OrderRequest) (*domain.Order, error) {
	var total float64
	var orderItems []domain.OrderItem

	// Validate products and calculate total
	for _, item := range req.Items {
		product, err := u.productRepo.GetByID(item.ProductID)
		if err != nil {
			return nil, errors.New("product not found")
		}

		if product.Stock < item.Quantity {
			return nil, errors.New("insufficient stock")
		}

		orderItem := domain.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		}

		total += product.Price * float64(item.Quantity)
		orderItems = append(orderItems, orderItem)
	}

	// Create order
	order := &domain.Order{
		UserID: userID,
		Total:  total,
		Status: "pending",
	}

	if err := u.orderRepo.Create(order); err != nil {
		return nil, err
	}

	// Create order items and update stock
	for _, item := range orderItems {
		item.OrderID = order.ID
		if err := u.orderRepo.CreateOrderItem(&item); err != nil {
			return nil, err
		}

		// Update product stock
		product, _ := u.productRepo.GetByID(item.ProductID)
		newStock := product.Stock - item.Quantity
		if err := u.productRepo.UpdateStock(item.ProductID, newStock); err != nil {
			return nil, err
		}
	}

	return order, nil
}

func (u *OrderUsecase) GetOrder(id uint) (*domain.Order, error) {
	order, err := u.orderRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	items, err := u.orderRepo.GetOrderItems(id)
	if err != nil {
		return nil, err
	}

	order.Items = make([]domain.OrderItem, len(items))
	for i, item := range items {
		order.Items[i] = *item
	}

	return order, nil
}

func (u *OrderUsecase) GetUserOrders(userID uint, page, limit int) ([]*domain.Order, error) {
	offset := (page - 1) * limit
	return u.orderRepo.GetByUserID(userID, limit, offset)
}
