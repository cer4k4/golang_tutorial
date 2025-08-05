package mysql

import (
    "database/sql"
    "shop/internal/domain"
)

type productRepository struct {
    db *sql.DB
}

func NewProductRepository(db *sql.DB) *productRepository {
    return &productRepository{db: db}
}

func (r *productRepository) Create(product *domain.Product) error {
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

func (r *productRepository) GetByID(id uint) (*domain.Product, error) {
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

func (r *productRepository) GetAll(limit, offset int) ([]*domain.Product, error) {
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

func (r *productRepository) Update(product *domain.Product) error {
    query := `UPDATE products SET name = ?, description = ?, price = ?, stock = ?, category = ? WHERE id = ?`
    _, err := r.db.Exec(query, product.Name, product.Description, product.Price, product.Stock, product.Category, product.ID)
    return err
}

func (r *productRepository) Delete(id uint) error {
    query := `DELETE FROM products WHERE id = ?`
    _, err := r.db.Exec(query, id)
    return err
}

func (r *productRepository) UpdateStock(id uint, stock int) error {
    query := `UPDATE products SET stock = ? WHERE id = ?`
    _, err := r.db.Exec(query, stock, id)
    return err
}