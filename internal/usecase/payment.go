package usecase

import (
	"shop/internal/domain"
	"shop/internal/repository"
)

type PaymentUsecase struct {
	paymentRepo repository.PaymentRepository
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
}

func NewPaymentUsecase(paymentRepo repository.PaymentRepository, orderRepo repository.OrderRepository, productRepo repository.ProductRepository) *PaymentUsecase {
	return &PaymentUsecase{
		paymentRepo: paymentRepo,
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (p *PaymentUsecase) CreatePayment(userID uint, req *domain.Payment) (*domain.Payment, error) {
	return &domain.Payment{}, nil
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
