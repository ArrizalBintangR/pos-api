# POS API - Point of Sale REST API

REST API backend untuk sistem Point of Sale dengan Golang dan PostgreSQL.

## Fitur

- JWT Authentication
- Role Based Access Control (RBAC) - 2 role: cashier, owner
- CRUD Sale Order
- CRUD User Cashier
- Pagination & Limit
- Standard Response Format

## Tech Stack

- Go 1.21+
- Gin Web Framework
- GORM (PostgreSQL)
- JWT (golang-jwt/jwt)
- bcrypt untuk password hashing

## Struktur Folder

```
├── config/          # Konfigurasi aplikasi
├── database/        # Database connection & migration
├── handlers/        # HTTP request handlers
├── middleware/      # Auth & RBAC middleware
├── models/          # Database models
├── routes/          # Route definitions
├── utils/           # Helper functions (JWT, response, pagination)
├── main.go          # Entry point
├── .env             # Environment variables
└── README.md
```

## Setup

### 1. Prerequisites

- Go 1.21 atau lebih baru
- PostgreSQL

### 2. Database Setup

```sql
CREATE DATABASE pos_db;
```

### 3. Environment Configuration

buat `.env`:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=pos_db
DB_SSL_MODE=disable

JWT_SECRET=your-super-secret-key-change-in-production
JWT_EXPIRY_HOURS=24

SERVER_PORT=8080

#GIN CONFIGURATION (USE debug FOR DEBUG MODE AND release FOR PRODUCTION MODE)
GIN_MODE=debug
```

### 4. Install Dependencies

```bash
go mod tidy
```

### 5. Run Application

```bash
go run main.go
```

Server akan berjalan di `http://localhost:8080`

## Default Users

Setelah pertama kali dijalankan, sistem akan membuat user default:

| Role | Username | Password |
|------|----------|----------|
| Owner | owner | owner123 |
| Cashier | cashier | cashier123 |

## API Endpoints

### Authentication

| Method | Endpoint | Description | Access |
|--------|----------|-------------|--------|
| POST | /auth/login | Login user | Public |
| POST | /auth/logout | Logout user | Authenticated |

### Sale Orders

| Method | Endpoint | Description | Access |
|--------|----------|-------------|--------|
| GET | /sale-orders | Get all sale orders (paginated) | Cashier, Owner |
| GET | /sale-orders/:id | Get sale order by ID | Cashier, Owner |
| POST | /sale-orders | Create sale order | Cashier, Owner |
| PATCH | /sale-orders/:id | Update sale order | Cashier, Owner |
| DELETE | /sale-orders/:id | Delete sale order | Cashier, Owner |

### User Cashier Management

| Method | Endpoint | Description | Access |
|--------|----------|-------------|--------|
| GET | /users/cashier | Get all cashiers (paginated) | Owner |
| GET | /users/cashier/:id | Get cashier by ID | Owner |
| POST | /users/cashier | Create cashier | Owner |
| PATCH | /users/cashier/:id | Update cashier | Owner |
| DELETE | /users/cashier/:id | Delete cashier | Owner |

## Response Format

### Success Response
```json
{
  "code": 200,
  "status": "success",
  "message": "successful",
  "data": {...}
}
```

### Error Response
```json
{
  "code": 400,
  "status": "failed",
  "message": "error message"
}
```

### Paginated Response
```json
{
  "code": 200,
  "status": "success",
  "message": "successful",
  "data": {
    "items": [...],
    "total_items": 100,
    "total_pages": 10,
    "page": 1,
    "limit": 10
  }
}
```

## API Examples

### Login
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "owner", "password": "owner123"}'
```

### Create Sale Order
```bash
curl -X POST http://localhost:8080/sale-orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "customer_name": "John Doe",
    "notes": "Rush order",
    "items": [
      {"product_name": "Product A", "quantity": 2, "unit_price": 10000},
      {"product_name": "Product B", "quantity": 1, "unit_price": 25000}
    ]
  }'
```

### Get Sale Orders with Pagination
```bash
curl -X GET "http://localhost:8080/sale-orders?page=1&limit=10" \
  -H "Authorization: Bearer <token>"
```

### Create Cashier (Owner only)
```bash
curl -X POST http://localhost:8080/users/cashier \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <owner_token>" \
  -d '{"username": "newcashier", "password": "password123", "name": "New Cashier"}'
```

## Testing RBAC

1. Login sebagai **cashier** - hanya bisa akses `/sale-orders/*`
2. Login sebagai **owner** - bisa akses `/sale-orders/*` dan `/users/cashier/*`
3. Cashier mencoba akses `/users/cashier` akan mendapat response 403 Forbidden

