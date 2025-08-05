package usecase

import (
	"errors"
	"shop/internal/domain"
	"shop/internal/repository"
)

type CartUsecase struct {
	cartRepo    repository.CartRepository
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
}

func NewCartUseCase(cartRepo *repository.CartRepository, orderRepo *repository.OrderRepository, productRepo *repository.ProductRepository) *CartUsecase {
	return &CartUsecase{cartRepo: *cartRepo, orderRepo: *orderRepo, productRepo: *productRepo}
}

func (u *CartUsecase) CreateCart(userID uint, req []*domain.CartItems) (*domain.Order, error) {
	var total float64
	var cartItems []domain.CartItems

	// Validate products and calculate total
	for _, item := range req {
		product, err := u.productRepo.GetByID(item.ProductId)
		if err != nil {
			return nil, errors.New("product not found")
		}

		if product.Stock < item.Quantity {
			return nil, errors.New("insufficient stock")
		}

		cartItem := domain.CartItems{
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
			Fee:       product.Price,
		}

		total += product.Price * float64(item.Quantity)
		cartItems = append(cartItems, cartItem)
	}

	// Create order
	cart := &domain.CartOrder{
		UserID: userID,
		Total:  total,
		Status: true,
	}

	if err := u.cartRepo.CreateCart(cart); err != nil {
		return nil, err
	}

	// Create order items and update stock
	for _, item := range cartItems {
		if err := u.cartRepo.CreateCartItems(&item); err != nil {
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

func (u *CartUsecase) GetCart(userid uint) (*domain.CartOrder, error) {
	cart, err := u.cartRepo.GetByUserID(userid)
	if err != nil {
		return nil, err
	}

	items, err := u.cartRepo.GetCartItems(cart.UserID)
	if err != nil {
		return nil, err
	}

	cart.Items = make([]domain.CartItems, len(items))
	for i, item := range items {
		cart.Items[i].ProductId = item.Items[0].ProductId
	}
	return cart, nil
}
