package mysql

import (
	"database/sql"
	"shop/internal/domain"
)

type paymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *paymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(payment *domain.Payment) error {
	query := `INSERT INTO payments (user_id, total, status) VALUES (?, ?, ?)`
	result, err := r.db.Exec(query, payment.UserID, payment.Total, payment.Status)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	payment.Id = uint(id)
	return nil
}

func (r *paymentRepository) GetByID(id uint) (*domain.Payment, error) {
	payment := &domain.Payment{}
	query := `SELECT id, user_id, total, status, created_at, updated_at FROM payments WHERE id = ?`

	err := r.db.QueryRow(query, id).Scan(
		&payment.Id, &payment.UserID, &payment.Total, &payment.Status,
		&payment.CreatedAt, &payment.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (r *paymentRepository) GetByUserID(userID uint, limit, offset int) ([]*domain.Payment, error) {
	query := `SELECT id, user_id, total, status, created_at, updated_at FROM payments WHERE user_id = ? LIMIT ? OFFSET ?`
	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []*domain.Payment
	for rows.Next() {
		payment := &domain.Payment{}
		err := rows.Scan(
			&payment.Id, &payment.UserID, &payment.Total, &payment.Status,
			&payment.CreatedAt, &payment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	return payments, nil
}

func (r *paymentRepository) Update(payment *domain.Payment) error {
	query := `UPDATE payments SET status = ? , order_id = ? WHERE id = ?`
	_, err := r.db.Exec(query, payment.Status, payment.OrderID, payment.Id)
	return err
}

func (r *paymentRepository) CreateOrderItem(item *domain.OrderItem) error {
	query := `INSERT INTO order_items (order_id, product_id, quantity, price) VALUES (?, ?, ?, ?)`
	result, err := r.db.Exec(query, item.OrderID, item.ProductID, item.Quantity, item.Price)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	item.ID = uint(id)
	return nil
}
