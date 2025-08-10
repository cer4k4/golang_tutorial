package mysql

import (
	"database/sql"
	"shop/internal/domain"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	query := `INSERT INTO users (username, email, password, role) VALUES (?, ?, ?, ?)`
	result, err := r.db.Exec(query, user.Username, user.Email, user.Password, user.Role)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = uint(id)
	return nil
}

func (r *userRepository) GetByID(id uint) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT id, username, total_cart, lock_cart, email, password, role, created_at, updated_at FROM users WHERE id = ?`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.TotalCart, &user.LockCart, &user.Email, &user.Password,
		&user.Role, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT id, username, total_cart, lock_cart, email, password, role, created_at, updated_at FROM users WHERE email = ?`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.TotalCart, &user.LockCart, &user.Email, &user.Password,
		&user.Role, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) Update(user *domain.User) error {
	query := `UPDATE users SET username = ?, email = ?, role = ? WHERE id = ?`
	_, err := r.db.Exec(query, user.Username, user.Email, user.Role, user.ID)
	return err
}

func (r *userRepository) Delete(id uint) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
