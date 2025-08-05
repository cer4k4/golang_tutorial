package usecase

import (
	"errors"
	"log"
	"shop/internal/domain"
	"shop/internal/repository"
	"time"
)

type CartUsecase struct {
	userRepo    repository.UserRepository
	cartRepo    repository.CartItemsRepository
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
}

func NewCartUseCase(cartRepo *repository.CartItemsRepository, orderRepo *repository.OrderRepository, productRepo *repository.ProductRepository, userRepo *repository.UserRepository) *CartUsecase {
	return &CartUsecase{cartRepo: *cartRepo, orderRepo: *orderRepo, productRepo: *productRepo, userRepo: *userRepo}
}

func (u *CartUsecase) CreateCart(userID uint, req *domain.RequestCart) error {
	oldCart, err := u.userRepo.GetByID(req.UserID)
	if err != nil {
		log.Println(err)
		return err
	}
	if oldCart.LockCart {
		return errors.New("You have Cart in payment processing")
	}

	// Validate products and calculate total
	for _, item := range req.Items {
		product, err := u.productRepo.GetByID(item.ProductId)
		if err != nil {
			return errors.New("product not found")
		}
		if product.Stock < item.Quantity {
			return errors.New("insufficient stock")
		}
		// TODO: Add or Delete Old product

		oldCart.TotalCart += product.Price * float64(item.Quantity)
	}

	// Create order items and update stock
	for _, item := range req.Items {
		item.CreatedAt = time.Now()
		item.UserId = oldCart.ID
		if err = u.cartRepo.CreateCartItems(req.UserID, &item); err != nil {
			return err
		}

		// Update product stock
		// product, _ := u.productRepo.GetByID(item.ProductId)
		// newStock := product.Stock - item.Quantity
		// if err := u.productRepo.UpdateStock(item.ProductId, newStock); err != nil {
		// 	return err
		// }
	}
	if err = u.userRepo.Update(oldCart); err != nil {
		return err
	}
	return nil
}

func (u *CartUsecase) GetCart(userid uint) (*domain.RequestCart, error) {
	usercart, err := u.userRepo.GetByID(userid)
	if err != nil {
		return nil, err
	}
	items, err := u.cartRepo.GetByUserID(usercart.ID)
	if err != nil {
		return nil, err
	}
	var respon domain.RequestCart
	respon.Items = *items
	respon.UserID = userid
	return &respon, nil
}
