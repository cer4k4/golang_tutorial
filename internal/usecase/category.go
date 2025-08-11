package usecase

import (
	"shop/internal/domain"
	"shop/internal/repository"
)

type CategoryUsecase struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryUsecase(categoryRepo repository.CategoryRepository) *CategoryUsecase {
	return &CategoryUsecase{categoryRepo: categoryRepo}
}

func (u *CategoryUsecase) CreateProduct(req *domain.Category) (*domain.Category, error) {
	category := &domain.Category{
		Name:        req.Name,
		Description: req.Description,
		FatherId:    req.FatherId,
	}

	if err := u.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (u *CategoryUsecase) GetProduct(id uint) (*domain.Category, error) {
	return u.categoryRepo.GetByID(id)
}

func (u *CategoryUsecase) GetProducts(page, limit int) ([]*domain.Category, error) {
	offset := (page - 1) * limit
	return u.categoryRepo.GetAll(limit, offset)
}

func (u *CategoryUsecase) UpdateCategory(id uint, req *domain.Category) (*domain.Category, error) {
	category, err := u.categoryRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	category.Name = req.Name
	category.Description = req.Description
	category.FatherId = req.FatherId

	if err := u.categoryRepo.Update(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (u *CategoryUsecase) DeleteProduct(id uint) error {
	return u.categoryRepo.Delete(id)
}
