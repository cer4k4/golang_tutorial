package usecase

import (
	"errors"
	"math/rand"
	"shop/internal/domain"
	"shop/internal/repository"
	"time"
)

type PaymentUsecase struct {
	paymentRepo repository.PaymentRepository
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
	userRepo    repository.UserRepository
	cartRepo    repository.CartItemsRepository
}

func NewPaymentUsecase(paymentRepo repository.PaymentRepository, orderRepo repository.OrderRepository, productRepo repository.ProductRepository, userRepo repository.UserRepository, cartRepo repository.CartItemsRepository) *PaymentUsecase {
	return &PaymentUsecase{
		paymentRepo: paymentRepo,
		orderRepo:   orderRepo,
		productRepo: productRepo,
		userRepo:    userRepo,
		cartRepo:    cartRepo,
	}
}

func (p *PaymentUsecase) CreatePayment(userID uint) (*domain.Payment, error) {
	user, err := p.userRepo.GetByID(userID)
	if user.LockCart {
		return &domain.Payment{}, errors.New("You have another payment in proccesing")
	}
	if user.TotalCart <= 0 {
		return &domain.Payment{}, errors.New("You're cart is empety")
	}
	user.LockCart = true
	p.userRepo.Update(user)
	var payment domain.Payment
	payment.UserID = userID
	payment.CreatedAt = time.Now()
	payment.Status = "paymentCreated"
	if err != nil {
		return &domain.Payment{}, err
	}
	payment.Total = user.TotalCart
	if err := p.paymentRepo.Create(&payment); err != nil {
		return &domain.Payment{}, err
	}
	if !randomPaymentFaild() {
		payment.Status = "faild"
		p.paymentRepo.Update(&payment)
		time.Sleep(time.Duration(15 * time.Second))
		return &domain.Payment{}, errors.New("Payment Is faild")
	}

	cartItems, err := p.cartRepo.GetByUserID(userID)
	if err != nil {
		return &domain.Payment{}, err
	}
	var order domain.Order
	order.CreatedAt = time.Now()
	order.Status = "pending"
	order.Total = user.TotalCart
	order.UserID = userID
	err = p.orderRepo.Create(&order)
	if err != nil {
		return &domain.Payment{}, err
	}
	for _, val := range cartItems {
		p.cartRepo.Delete(&domain.CartItems{UserId: userID, ProductId: val.ProductId})
		p.orderRepo.CreateOrderItem(&domain.OrderItem{OrderID: order.ID, ProductID: val.ProductId, Quantity: val.Quantity, Price: val.Fee})
	}
	payment.Status = "paid"
	p.paymentRepo.Update(&payment)
	user.LockCart = false
	user.TotalCart = 0
	p.userRepo.Update(user)
	time.Sleep(time.Duration(10 * time.Second))
	return &payment, err
}

func (p *PaymentUsecase) GetPayment(id uint) (*domain.Payment, error) {
	payment, err := p.paymentRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (p *PaymentUsecase) GetUserPayments(userID uint, page, limit int) ([]*domain.Payment, error) {
	offset := (page - 1) * limit
	return p.paymentRepo.GetByUserID(userID, limit, offset)
}

func randomPaymentFaild() bool {
	rand.Seed(time.Now().UnixNano())
	// Generate a random boolean
	randomBool := rand.Float64() < 0.5
	return randomBool
}
