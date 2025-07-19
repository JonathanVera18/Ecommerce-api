# Environment Configuration Guide

## Quick Setup

1. **Copy the example environment file:**
   ```bash
   cp .env.example .env
   ```

2. **Update the key values in `.env`:**
   - Database credentials
   - JWT secret (make it long and secure)
   - Email SMTP settings
   - Stripe API keys (if using payments)

## Important Configuration Notes

### 🔒 Security Settings

**JWT_SECRET**: This should be a long, random string in production. Generate one using:
```bash
openssl rand -base64 64
```

**Database Password**: Use a strong password, especially in production.

**BCRYPT_COST**: Higher values (12-15) are more secure but slower. 12 is recommended for production.

### 📧 Email Configuration

For **Gmail SMTP**:
1. Enable 2-factor authentication
2. Generate an "App Password" (not your regular password)
3. Use the app password in `SMTP_PASSWORD`

For **Production**: Consider using dedicated email services like:
- SendGrid
- AWS SES
- Mailgun
- Postmark

### 💳 Payment Configuration

**Development**: Use Stripe test keys (start with `sk_test_` and `pk_test_`)
**Production**: Use live keys (start with `sk_live_` and `pk_live_`)

### 🗄️ Database Configuration

**Development**: 
- Port 5433 is used to avoid conflicts with local PostgreSQL (default 5432)
- Docker Compose automatically sets up PostgreSQL on port 5433

**Production**:
- Use managed database services (AWS RDS, Google Cloud SQL, etc.)
- Enable SSL mode (`DB_SSL_MODE=require`)
- Use connection pooling

### 📁 File Upload Configuration

**MAX_FILE_SIZE**: Size in bytes (10485760 = 10MB)
**ALLOWED_FILE_TYPES**: Comma-separated list of allowed extensions
**UPLOAD_DIR**: Local directory for file storage (consider cloud storage for production)

### 🚀 Performance Settings

**Redis Configuration**:
- `REDIS_MAX_IDLE`: Number of idle connections in pool
- `REDIS_MAX_ACTIVE`: Maximum active connections

**Cache TTL Settings**:
- `CACHE_TTL`: General cache duration
- `PRODUCT_CACHE_TTL`: Product-specific cache (shorter due to stock changes)
- `USER_CACHE_TTL`: User data cache

### 🔍 Feature Flags

Enable/disable features using boolean flags:
- `ENABLE_REGISTRATION`: Allow new user sign-ups
- `ENABLE_REVIEWS`: Product review system
- `ENABLE_WISHLIST`: User wishlist functionality
- `ENABLE_CART`: Shopping cart system
- `ENABLE_NOTIFICATIONS`: Push notifications
- `ENABLE_FILE_UPLOAD`: File upload endpoints

### 🌍 Environment-Specific Settings

#### Development
```env
APP_ENV=development
DEBUG=true
MOCK_PAYMENTS=true
MOCK_EMAILS=true
LOG_LEVEL=debug
```

#### Staging
```env
APP_ENV=staging
DEBUG=false
MOCK_PAYMENTS=true
MOCK_EMAILS=false
LOG_LEVEL=info
```

#### Production
```env
APP_ENV=production
DEBUG=false
MOCK_PAYMENTS=false
MOCK_EMAILS=false
LOG_LEVEL=warn
```

## Docker Configuration

When using Docker Compose, the following values are automatically configured:
- `DB_HOST=postgres` (container name)
- `REDIS_HOST=redis` (container name)
- Database credentials match docker-compose.yml

## Environment Variables Priority

1. System environment variables (highest priority)
2. `.env` file
3. Default values in code (lowest priority)

## Security Checklist for Production

- [ ] Change all default passwords
- [ ] Use strong JWT secret (64+ characters)
- [ ] Enable database SSL (`DB_SSL_MODE=require`)
- [ ] Use HTTPS for all URLs
- [ ] Set secure CORS origins
- [ ] Use environment-specific secrets
- [ ] Enable rate limiting
- [ ] Set appropriate log levels
- [ ] Disable debug mode
- [ ] Use real email service (not mock)
- [ ] Use live payment keys
- [ ] Set up monitoring and health checks

## Common Issues

**Database Connection Failed**:
- Check if PostgreSQL is running
- Verify port (5433 for Docker, 5432 for local)
- Check credentials and database name

**Redis Connection Failed**:
- Ensure Redis is running
- Check port and host settings
- Verify Redis configuration

**Email Not Sending**:
- Check SMTP credentials
- Verify app password for Gmail
- Test with a simple SMTP client first

**File Upload Issues**:
- Check upload directory permissions
- Verify max file size settings
- Ensure allowed file types are correct

## Monitoring Variables

Add these for production monitoring:
```env
# New Relic
NEW_RELIC_LICENSE_KEY=your_license_key
NEW_RELIC_APP_NAME=ecommerce-api

# Sentry (Error tracking)
SENTRY_DSN=https://your_sentry_dsn

# DataDog
DD_API_KEY=your_datadog_api_key
```
