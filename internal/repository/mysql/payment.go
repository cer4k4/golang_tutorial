package mysql

import (
	"database/sql"
	"shop/internal/domain"
	"time"
)

type paymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *paymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(payment *domain.Payment) error {
	query := `INSERT INTO payments (user_id, amount, status, payment_method, gateway_transaction_id, gateway_response, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := r.db.Exec(query,
		payment.UserID,
		payment.Amount,
		payment.Status,
		payment.PaymentMethod,
		payment.GatewayTransactionID,
		payment.GatewayResponse,
		now,
		now,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	payment.ID = uint(id)
	payment.CreatedAt = now
	payment.UpdatedAt = now

	return nil
}

func (r *paymentRepository) GetByID(id uint) (*domain.Payment, error) {
	query := `SELECT id, user_id, amount, status, payment_method, gateway_transaction_id, gateway_response, created_at, updated_at 
	          FROM payments WHERE id = ?`

	payment := &domain.Payment{}
	err := r.db.QueryRow(query, id).Scan(
		&payment.ID,
		&payment.UserID,
		&payment.Amount,
		&payment.Status,
		&payment.PaymentMethod,
		&payment.GatewayTransactionID,
		&payment.GatewayResponse,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (r *paymentRepository) GetByUserID(userID uint, limit, offset int) ([]*domain.Payment, error) {
	query := `SELECT id, user_id, amount, status, payment_method, gateway_transaction_id, gateway_response, created_at, updated_at 
	          FROM payments WHERE user_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []*domain.Payment
	for rows.Next() {
		payment := &domain.Payment{}
		err := rows.Scan(
			&payment.ID,
			&payment.UserID,
			&payment.Amount,
			&payment.Status,
			&payment.PaymentMethod,
			&payment.GatewayTransactionID,
			&payment.GatewayResponse,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	return payments, nil
}

func (r *paymentRepository) Update(payment *domain.Payment) error {
	query := `UPDATE payments SET amount = ?, status = ?, payment_method = ?, gateway_transaction_id = ?, gateway_response = ?, updated_at = ? 
	          WHERE id = ?`

	payment.UpdatedAt = time.Now()
	_, err := r.db.Exec(query,
		payment.Amount,
		payment.Status,
		payment.PaymentMethod,
		payment.GatewayTransactionID,
		payment.GatewayResponse,
		payment.UpdatedAt,
		payment.ID,
	)

	return err
}

func (r *paymentRepository) GetPendingPaymentByUserID(userID uint) (*domain.Payment, error) {
	query := `SELECT id, user_id, amount, status, payment_method, gateway_transaction_id, gateway_response, created_at, updated_at 
	          FROM payments WHERE user_id = ? AND status = 'pending' ORDER BY created_at DESC LIMIT 1`

	payment := &domain.Payment{}
	err := r.db.QueryRow(query, userID).Scan(
		&payment.ID,
		&payment.UserID,
		&payment.Amount,
		&payment.Status,
		&payment.PaymentMethod,
		&payment.GatewayTransactionID,
		&payment.GatewayResponse,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return payment, nil
}
