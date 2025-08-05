package mysql

import (
    "database/sql"
    "shop/internal/domain"
)

type orderRepository struct {
    db *sql.DB
}

func NewOrderRepository(db *sql.DB) *orderRepository {
    return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *domain.Order) error {
    query := `INSERT INTO orders (user_id, total, status) VALUES (?, ?, ?)`
    result, err := r.db.Exec(query, order.UserID, order.Total, order.Status)
    if err != nil {
        return err
    }
    
    id, err := result.LastInsertId()
    if err != nil {
        return err
    }
    
    order.ID = uint(id)
    return nil
}

func (r *orderRepository) GetByID(id uint) (*domain.Order, error) {
    order := &domain.Order{}
    query := `SELECT id, user_id, total, status, created_at, updated_at FROM orders WHERE id = ?`
    
    err := r.db.QueryRow(query, id).Scan(
        &order.ID, &order.UserID, &order.Total, &order.Status,
        &order.CreatedAt, &order.UpdatedAt,
    )
    
    if err != nil {
        return nil, err
    }
    
    return order, nil
}

func (r *orderRepository) GetByUserID(userID uint, limit, offset int) ([]*domain.Order, error) {
    query := `SELECT id, user_id, total, status, created_at, updated_at FROM orders WHERE user_id = ? LIMIT ? OFFSET ?`
    rows, err := r.db.Query(query, userID, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var orders []*domain.Order
    for rows.Next() {
        order := &domain.Order{}
        err := rows.Scan(
            &order.ID, &order.UserID, &order.Total, &order.Status,
            &order.CreatedAt, &order.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        orders = append(orders, order)
    }

    return orders, nil
}

func (r *orderRepository) Update(order *domain.Order) error {
    query := `UPDATE orders SET status = ?, total = ? WHERE id = ?`
    _, err := r.db.Exec(query, order.Status, order.Total, order.ID)
    return err
}

func (r *orderRepository) CreateOrderItem(item *domain.OrderItem) error {
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

func (r *orderRepository) GetOrderItems(orderID uint) ([]*domain.OrderItem, error) {
    query := `
        SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, oi.price,
               p.id, p.name, p.description, p.price, p.stock, p.category, p.created_at, p.updated_at
        FROM order_items oi
        JOIN products p ON oi.product_id = p.id
        WHERE oi.order_id = ?`
    
    rows, err := r.db.Query(query, orderID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []*domain.OrderItem
    for rows.Next() {
        item := &domain.OrderItem{}
        err := rows.Scan(
            &item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price,
            &item.Product.ID, &item.Product.Name, &item.Product.Description,
            &item.Product.Price, &item.Product.Stock, &item.Product.Category,
            &item.Product.CreatedAt, &item.Product.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        items = append(items, item)
    }

    return items, nil
}