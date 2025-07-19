# E-commerce API

A comprehensive e-commerce backend API built with Go, featuring user management, product catalog, order processing, reviews, and payment integration.

## Features

### Core Features
- **User Management**: Registration, authentication, profile management with role-based access (Customer, Seller, Admin)
- **Product Management**: Full CRUD operations with inventory tracking, categories, and image management
- **Order Processing**: Shopping cart, checkout, payment processing, and order status tracking
- **Review System**: Product reviews and ratings with helpful votes and seller responses
- **Email Notifications**: Automated emails for order confirmations, shipping updates, and more
- **Analytics**: Sales reports and business intelligence for sellers and admins

### Technical Features
- **RESTful API**: Clean REST endpoints with proper HTTP status codes
- **JWT Authentication**: Secure token-based authentication
- **Role-based Authorization**: Different access levels for customers, sellers, and admins
- **Database Migrations**: Version-controlled database schema changes
- **Payment Integration**: Stripe payment processing
- **Email Service**: SMTP email notifications
- **Docker Support**: Containerized application with Docker Compose
- **Input Validation**: Comprehensive request validation
- **Error Handling**: Structured error responses
- **Pagination**: Efficient data pagination for large datasets

## Tech Stack

- **Language**: Go 1.21+
- **Web Framework**: Echo v4
- **Database**: PostgreSQL with GORM ORM
- **Cache**: Redis
- **Authentication**: JWT tokens
- **Payment**: Stripe
- **Email**: SMTP
- **Documentation**: Swagger/OpenAPI
- **Containerization**: Docker & Docker Compose

## Project Structure

```
ecommerce-api/
├── internal/                   # Private application code
│   ├── config/                # Configuration management
│   ├── models/                # Data models
│   ├── repository/            # Data access layer
│   ├── service/               # Business logic layer
│   ├── handler/               # HTTP handlers
│   ├── middleware/            # HTTP middleware
│   └── utils/                 # Utility functions
├── pkg/                       # Public packages
│   ├── email/                 # Email service
│   └── payment/               # Payment service
├── migrations/                # Database migrations
├── scripts/                   # Utility scripts
├── test/                      # Test files
├── docs/                      # Documentation
├── api/                       # API documentation
├── main.go                    # Application entry point
├── docker-compose.yml         # Docker composition
├── Dockerfile                 # Container definition
├── Makefile                   # Build automation
└── README.md                  # Project documentation
```

## Quick Start

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 12+
- Redis 6+
- Docker & Docker Compose (optional)

### Environment Setup

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd ecommerce-api
   ```

2. **Copy environment file**:
   ```bash
   copy .env.example .env
   ```

3. **Update environment variables** in `.env`:
   ```env
   # Database
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=your_password
   DB_NAME=ecommerce_db
   
   # JWT Secret
   JWT_SECRET=your-super-secret-jwt-key
   
   # Stripe Keys
   STRIPE_SECRET_KEY=sk_test_...
   STRIPE_PUBLISHABLE_KEY=pk_test_...
   
   # Email Configuration
   SMTP_HOST=smtp.gmail.com
   SMTP_PORT=587
   SMTP_USERNAME=your-email@gmail.com
   SMTP_PASSWORD=your-app-password
   ```

### Running with Docker (Recommended)

1. **Start all services**:
   ```bash
   docker-compose up -d
   ```

2. **Check logs**:
   ```bash
   docker-compose logs -f api
   ```

3. **Stop services**:
   ```bash
   docker-compose down
   ```

### Running Locally

1. **Install dependencies**:
   ```bash
   go mod download
   ```

2. **Set up database**:
   ```bash
   # Create database
   createdb ecommerce_db
   
   # Run migrations
   make migrate-up
   ```

3. **Start the application**:
   ```bash
   make run
   ```

4. **Create admin user** (optional):
   ```bash
   make create-admin
   ```

## API Documentation

### Authentication Endpoints

- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh JWT token
- `POST /api/v1/auth/logout` - User logout
- `POST /api/v1/auth/change-password` - Change password

### User Endpoints

- `GET /api/v1/users/profile` - Get user profile
- `PUT /api/v1/users/profile` - Update user profile
- `GET /api/v1/users` - List users (Admin only)
- `POST /api/v1/users` - Create user (Admin only)
- `GET /api/v1/users/{id}` - Get user by ID (Admin only)
- `PUT /api/v1/users/{id}` - Update user (Admin only)
- `DELETE /api/v1/users/{id}` - Delete user (Admin only)

### Product Endpoints

- `GET /api/v1/products` - List products
- `GET /api/v1/products/{id}` - Get product by ID
- `GET /api/v1/products/slug/{slug}` - Get product by slug
- `POST /api/v1/products` - Create product (Seller/Admin)
- `PUT /api/v1/products/{id}` - Update product (Seller/Admin)
- `DELETE /api/v1/products/{id}` - Delete product (Seller/Admin)
- `GET /api/v1/products/search` - Search products
- `GET /api/v1/products/category/{category}` - Get products by category
- `GET /api/v1/products/featured` - Get featured products

### Order Endpoints

- `GET /api/v1/orders` - List orders
- `GET /api/v1/orders/{id}` - Get order by ID
- `POST /api/v1/orders` - Create order
- `PUT /api/v1/orders/{id}/status` - Update order status
- `POST /api/v1/orders/{id}/cancel` - Cancel order
- `POST /api/v1/orders/payment` - Process payment

### Cart Endpoints

- `GET /api/v1/cart` - Get cart
- `POST /api/v1/cart/items` - Add item to cart
- `PUT /api/v1/cart/items` - Update cart item
- `DELETE /api/v1/cart/items/{productId}` - Remove item from cart
- `DELETE /api/v1/cart` - Clear cart

### Review Endpoints

- `GET /api/v1/reviews` - List reviews
- `GET /api/v1/reviews/{id}` - Get review by ID
- `POST /api/v1/reviews` - Create review
- `PUT /api/v1/reviews/{id}` - Update review
- `DELETE /api/v1/reviews/{id}` - Delete review
- `POST /api/v1/reviews/{id}/helpful` - Mark review as helpful
- `POST /api/v1/reviews/{id}/response` - Add seller response

### Admin Endpoints

- `GET /api/v1/admin/stats/users` - User statistics
- `GET /api/v1/admin/stats/products` - Product statistics
- `GET /api/v1/admin/stats/orders` - Order statistics
- `GET /api/v1/admin/stats/reviews` - Review statistics

## Database Schema

The database consists of the following main tables:

- **users**: User accounts and profiles
- **products**: Product catalog
- **product_images**: Product image management
- **orders**: Customer orders
- **order_items**: Items within orders
- **carts**: Shopping carts
- **cart_items**: Items in shopping carts
- **reviews**: Product reviews and ratings
- **review_helpful**: Helpful votes on reviews

## Development

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
```

### Code Quality

```bash
# Format code
make fmt

# Lint code
make lint

# Security scan
make security
```

### Database Operations

```bash
# Run migrations
make migrate-up

# Rollback migrations
make migrate-down

# Seed database
make seed
```

## Deployment

### Production Environment

1. **Set production environment variables**
2. **Build Docker image**:
   ```bash
   docker build -t ecommerce-api:latest .
   ```

3. **Deploy with Docker Compose**:
   ```bash
   docker-compose -f docker-compose.prod.yml up -d
   ```

### Health Checks

The API provides a health check endpoint:
- `GET /health` - Returns application health status

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | `password` |
| `DB_NAME` | Database name | `ecommerce_db` |
| `REDIS_HOST` | Redis host | `localhost` |
| `REDIS_PORT` | Redis port | `6379` |
| `JWT_SECRET` | JWT signing secret | Required |
| `SERVER_PORT` | Server port | `8080` |
| `STRIPE_SECRET_KEY` | Stripe secret key | Required |
| `SMTP_HOST` | SMTP host | Required |
| `SMTP_USERNAME` | SMTP username | Required |
| `SMTP_PASSWORD` | SMTP password | Required |

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for your changes
5. Run the test suite
6. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For support and questions:
- Create an issue on GitHub
- Email: support@ecommerce.com
- Documentation: [API Docs](./docs/API.md)

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for a list of changes and version history.


