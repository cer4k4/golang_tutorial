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

func (r *cartItemRepository) GetByUserID(userID uint) ([]*domain.CartItems, error) {
	query := `SELECT user_id, total, status, created_at, updated_at FROM orders WHERE user_id = ?`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.CartItems
	for rows.Next() {
		item := &domain.CartItems{}
		err := rows.Scan(
			&item.ProductId,
			&item.Fee,
			&item.Quantity,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, item)
	}
	return orders, nil
}

func (r *cartItemRepository) Update(cart *domain.CartItems) error {
	cart.UpdatedAt = time.Now()
	query := `UPDATE orders SET quantity = ?, fee = ? WHERE product_id = ? AND user_id = ?`
	_, err := r.db.Exec(query, cart.Quantity, cart.Fee, cart.ProductId, cart.UserId)
	return err
}

func (r *cartItemRepository) CreateCartItems(userid uint, item *domain.CartItems) error {
	query := `INSERT INTO cart_items (user_id, product_id, quantity, fee) VALUES (?, ?, ?, ?)`
	_, err := r.db.Exec(query, userid, item.ProductId, item.Quantity, item.Fee)
	if err != nil {
		return err
	}
	return nil
}

func (r *cartItemRepository) GetCartItems(orderID uint) ([]*domain.CartItems, error) {
	query := `
        SELECT oi.id, oi.user_id, oi.product_id, oi.quantity, oi.price,
               p.id, p.name, p.description, p.price, p.stock, p.category, p.created_at, p.updated_at
        FROM cart_items oi
        JOIN products p ON oi.product_id = p.id
        WHERE oi.order_id = ?`

	rows, err := r.db.Query(query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*domain.CartItems
	for rows.Next() {
		item := &domain.CartItems{}
		err := rows.Scan(&item.ProductId, &item.Quantity, &item.Fee, &item.Quantity)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
