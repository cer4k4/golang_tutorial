package mysql

import (
	"database/sql"
	"shop/internal/domain"
)

type cacategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(product *domain.Product) error {
	query := `INSERT INTO products (name, description, price, stock, category) VALUES (?, ?, ?, ?, ?)`
	result, err := r.db.Exec(query, product.Name, product.Description, product.Price, product.Stock, product.Category)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	product.ID = uint(id)
	return nil
}

func (r *categoryRepository) GetByID(id uint) (*domain.Product, error) {
	product := &domain.Product{}
	query := `SELECT id, name, description, price, stock, category, created_at, updated_at FROM products WHERE id = ?`

	err := r.db.QueryRow(query, id).Scan(
		&product.ID, &product.Name, &product.Description, &product.Price,
		&product.Stock, &product.Category, &product.CreatedAt, &product.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *categoryRepository) GetAll(limit, offset int) ([]*domain.Product, error) {
	query := `SELECT id, name, description, price, stock, category, created_at, updated_at FROM products LIMIT ? OFFSET ?`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		product := &domain.Product{}
		err := rows.Scan(
			&product.ID, &product.Name, &product.Description, &product.Price,
			&product.Stock, &product.Category, &product.CreatedAt, &product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (r *categoryRepository) Update(category *domain.Category) error {
	query := `UPDATE products SET name = ?, description = ?, father_id = ? WHERE id = ?`
	_, err := r.db.Exec(query, category.Name, category.Description, category.FatherId, cacategory.ID)
	return err
}

func (r *categoryRepository) Delete(id uint) error {
	query := `DELETE FROM products WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *categoryRepository) UpdateStock(id uint, stock int) error {
	query := `UPDATE products SET stock = ? WHERE id = ?`
	_, err := r.db.Exec(query, stock, id)
	return err
}
