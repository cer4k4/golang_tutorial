package mysql

import (
	"database/sql"
	"shop/internal/domain"
	"time"
)

type cartItemRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) *cartItemRepository {
	return &cartItemRepository{db: db}
}

func (r *cartItemRepository) GetByUserID(userID uint) ([]domain.CartItems, error) {
	query := `SELECT product_id, quantity, fee, user_id, created_at, updated_at 
	          FROM cart_items WHERE user_id = ?`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cartItems []domain.CartItems
	for rows.Next() {
		item := domain.CartItems{}
		err := rows.Scan(
			&item.ProductId,
			&item.Quantity,
			&item.Fee,
			&item.UserId,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		cartItems = append(cartItems, item)
	}
	return cartItems, nil
}

func (r *cartItemRepository) CreateCartItems(userID uint, item domain.CartItems) error {
	query := `INSERT INTO cart_items (user_id, product_id, quantity, fee, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?)`
	now := time.Now()
	_, err := r.db.Exec(query, userID, item.ProductId, item.Quantity, item.Fee, now, now)
	return err
}

func (r *cartItemRepository) UpdateCartItem(userID uint, item domain.CartItems) error {
	query := `UPDATE cart_items SET quantity = ?, fee = ?, updated_at = ? 
	          WHERE product_id = ? AND user_id = ?`
	_, err := r.db.Exec(query, item.Quantity, item.Fee, time.Now(), item.ProductId, userID)
	return err
}

func (r *cartItemRepository) DeleteCartItem(userID, productID uint) error {
	query := `DELETE FROM cart_items WHERE user_id = ? AND product_id = ?`
	_, err := r.db.Exec(query, userID, productID)
	return err
}

func (r *cartItemRepository) ClearCart(userID uint) error {
	query := `DELETE FROM cart_items WHERE user_id = ?`
	_, err := r.db.Exec(query, userID)
	return err
}

func (r *cartItemRepository) GetCartItemByUserAndProduct(userID, productID uint) (*domain.CartItems, error) {
	query := `SELECT product_id, quantity, fee, user_id, created_at, updated_at 
	          FROM cart_items WHERE user_id = ? AND product_id = ?`

	item := &domain.CartItems{}
	err := r.db.QueryRow(query, userID, productID).Scan(
		&item.ProductId,
		&item.Quantity,
		&item.Fee,
		&item.UserId,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return item, nil
}
