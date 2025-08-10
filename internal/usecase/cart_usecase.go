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

// TODO: add qauntiy from stock
func (u *CartUsecase) CreateCart(userID uint, req *domain.RequestCart) error {
	oldCart, err := u.userRepo.GetByID(userID)
	if err != nil {
		log.Println(err)
		return err
	}
	if oldCart.LockCart {
		return errors.New("You have Cart in payment processing")
	}
	type newModel struct {
		indexlist int
		s         domain.CartItems
	}
	var oldcart = make(map[uint]newModel)
	oldItems, err := u.GetCart(userID)
	for o := range oldItems.Items {
		var m newModel
		m.s = oldItems.Items[o]
		m.indexlist = o
		oldcart[oldItems.Items[o].ProductId] = m
	}
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
		// if that item exist
		if len(oldItems.Items) != 0 {
			if oldcart[item.ProductId].s.ProductId == item.ProductId {
				oldCart.TotalCart += product.Price * float64(item.Quantity)
				oldItems.Items[oldcart[item.ProductId].indexlist].Quantity += item.Quantity
				if oldItems.Items[oldcart[item.ProductId].indexlist].Quantity < 0 {
					u.cartRepo.Delete(&oldItems.Items[oldcart[item.ProductId].indexlist])
				} else {
					if err = u.cartRepo.Update(&oldItems.Items[oldcart[item.ProductId].indexlist]); err != nil {
						return err
					}
				}
			} else {
				item.CreatedAt = time.Now()
				item.UserId = userID
				item.Fee = product.Price
				if item.Quantity > 0 {
					oldCart.TotalCart += product.Price * float64(item.Quantity)

					if err = u.cartRepo.CreateCartItems(userID, &item); err != nil {
						return err
					}
				}
			}
		} else {
			item.CreatedAt = time.Now()
			item.UserId = userID
			item.Fee = product.Price
			if item.Quantity > 0 {
				if err = u.cartRepo.CreateCartItems(userID, &item); err != nil {
					return err
				}
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
