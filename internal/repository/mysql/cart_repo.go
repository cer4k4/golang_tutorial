package mysql

import (
	"database/sql"
	"shop/internal/domain"
)

type cartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) *cartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) CreateCart(cart *domain.CartOrder) error {
	query := `INSERT INTO carts (user_id, total, status) VALUES (?, ?, ?)`
	_, err := r.db.Exec(query, cart.UserID, cart.Total, cart.Status)
	if err != nil {
		return err
	}
	return nil
}

func (r *cartRepository) GetByUserID(userID uint) ([]*domain.CartOrder, error) {
	query := `SELECT user_id, total, status, created_at, updated_at FROM orders WHERE user_id = ?`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.CartOrder
	for rows.Next() {
		order := &domain.CartOrder{}
		err := rows.Scan(
			&order.UserID,
			&order.Total,
			&order.Status,
			&order.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *cartRepository) Update(order *domain.CartOrder) error {
	query := `UPDATE orders SET status = ?, total = ? WHERE id = ?`
	_, err := r.db.Exec(query, order.Status, order.Total)
	return err
}

func (r *cartRepository) CreateCartItems(userid uint, item *domain.CartItems) error {
	query := `INSERT INTO cart_items (user_id, product_id, quantity, fee) VALUES (?, ?, ?, ?)`
	_, err := r.db.Exec(query, userid, item.ProductId, item.Quantity, item.Fee)
	if err != nil {
		return err
	}
	return nil
}

func (r *cartRepository) GetCartItems(orderID uint) ([]*domain.CartItems, error) {
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
