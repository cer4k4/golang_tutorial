package usecase

import (
	"shop/internal/domain"
	"shop/internal/repository"
)

type ProductUsecase struct {
	productRepo repository.ProductRepository
}

func NewProductUsecase(productRepo repository.ProductRepository) *ProductUsecase {
	return &ProductUsecase{productRepo: productRepo}
}

func (u *ProductUsecase) CreateProduct(req *domain.ProductRequest) (*domain.Product, error) {
	product := &domain.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Discount:    req.Discount,
		Category:    req.Category,
	}

	if err := u.productRepo.Create(product); err != nil {
		return nil, err
	}

	return product, nil
}

func (u *ProductUsecase) GetProduct(id uint) (*domain.Product, error) {
	return u.productRepo.GetByID(id)
}

func (u *ProductUsecase) GetProducts(page, limit int) ([]*domain.Product, error) {
	offset := (page - 1) * limit
	return u.productRepo.GetAll(limit, offset)
}

func (u *ProductUsecase) UpdateProduct(id uint, req *domain.ProductRequest) (*domain.Product, error) {
	product, err := u.productRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Stock = req.Stock
	product.Discount = req.Discount
	product.Category = req.Category

	if err := u.productRepo.Update(product); err != nil {
		return nil, err
	}

	return product, nil
}

func (u *ProductUsecase) DeleteProduct(id uint) error {
	return u.productRepo.Delete(id)
}
