# Go Shop - Monolithic E-Commerce API

A monolithic e-commerce API built with Go, following Clean Architecture principles.

## Features

- User authentication with JWT tokens
- Product management (CRUD operations)
- Order management
- Role-based access control (User/Admin)
- Clean Architecture implementation
- MySQL database integration

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin
- **Database**: MySQL
- **Authentication**: JWT
- **Architecture**: Clean Architecture

## Project Structure

```
/
├── cmd/                    # Application entry points
├── internal/               # Private application code
│   ├── domain/            # Business entities
│   ├── repository/        # Data access layer
│   ├── usecase/          # Business logic
│   ├── delivery/         # Delivery mechanisms (HTTP)
│   └── middleware/       # HTTP middleware
├── pkg/                   # Public packages
└── go.mod
```

## Setup

1. **Clone the repository**
2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Setup MySQL database**:
   - Create a database named `shop_db`
   - Update connection details in environment variables

4. **Environment Variables**:
   ```bash
   # Database
   DB_HOST=localhost
   DB_PORT=3306
   DB_USER=root
   DB_PASSWORD=your_password
   DB_NAME=shop_db
   
   # JWT
   JWT_SECRET=your-super-secret-jwt-key
   ```

5. **Run the application**:
   ```bash
   go run cmd/main.go
   ```

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login user

### Users
- `GET /api/v1/users/profile` - Get user profile (authenticated)

### Products
- `GET /api/v1/products` - Get all products
- `GET /api/v1/products/:id` - Get product by ID
- `POST /api/v1/admin/products` - Create product (admin only)
- `PUT /api/v1/admin/products/:id` - Update product (admin only)
- `DELETE /api/v1/admin/products/:id` - Delete product (admin only)

### Orders
- `POST /api/v1/orders` - Create order (authenticated)
- `GET /api/v1/orders` - Get user orders (authenticated)
- `GET /api/v1/orders/:id` - Get order by ID (authenticated)

## Example Usage

### Register a new user
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Create a product (Admin only)
```bash
curl -X POST http://localhost:8080/api/v1/admin/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "iPhone 15",
    "description": "Latest iPhone model",
    "price": 999.99,
    "stock": 100,
    "category": "Electronics"
  }'
```

### Get products
```bash
curl http://localhost:8080/api/v1/products?page=1&limit=10
```

### Create an order
```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "items": [
      {
        "product_id": 1,
        "quantity": 2
      }
    ]
  }'
```

## Database Schema

The application automatically creates the following tables:

- `users` - User information and authentication
- `products` - Product catalog
- `orders` - Order information
- `order_items` - Order line items

## Development

To run in development mode:

```bash
# Install air for hot reloading (optional)
go install github.com/cosmtrek/air@latest

# Run with hot reloading
air
```

## Testing

To create an admin user, you can directly insert into the database:

```sql
INSERT INTO users (username, email, password, role) 
VALUES ('admin', 'admin@example.com', '$2a$10$hashed_password', 'admin');
```

Or register a user and manually update the role in the database.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License.

// .env.example
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=shop_db

# JWT Configuration  
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# Server Configuration
PORT=8080