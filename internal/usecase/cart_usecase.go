package usecase

import (
	"errors"
	"log"
	"shop/internal/domain"
	"shop/internal/repository"
	"time"
)

type CartUsecase struct {
	cartRepo    repository.CartItemsRepository
	userRepo    repository.UserRepository
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
}

func NewCartUseCase(cartRepo repository.CartItemsRepository, userRepo repository.UserRepository, orderRepo repository.OrderRepository, productRepo repository.ProductRepository) *CartUsecase {
	return &CartUsecase{cartRepo: cartRepo, userRepo: userRepo, orderRepo: orderRepo, productRepo: productRepo}
}

func (u *CartUsecase) CreateCart(userID uint, req *domain.RequestCart) error {
	oldCart, err := u.userRepo.GetByID(userID)
	if err != nil {
		log.Println(err)
		return err
	}
	if oldCart.LockCart {
		return errors.New("You have Cart in payment processing")
	}
	oldItems, err := u.GetCart(userID)
	if err != nil {
		return err
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

		// if that item exist
		if len(oldItems.Items) != 0 {
			for o := range oldItems.Items {
				if oldItems.Items[o].ProductId == item.ProductId {
					oldItems.Items[o].Quantity += item.Quantity
					if oldItems.Items[o].Quantity == 0 {
						u.cartRepo.Delete(&oldItems.Items[o])
					} else {
						if err = u.cartRepo.Update(&oldItems.Items[o]); err != nil {
							return err
						}
					}
				} else {
					item.CreatedAt = time.Now()
					item.UserId = userID
					item.Fee = product.Price
					if err = u.cartRepo.CreateCartItems(userID, &item); err != nil {
						return err
					}
				}
			}
		} else {
			item.CreatedAt = time.Now()
			item.UserId = userID
			item.Fee = product.Price
			if err = u.cartRepo.CreateCartItems(userID, &item); err != nil {
				return err
			}
		}
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
	respon.Items = items
	return &respon, nil
}
