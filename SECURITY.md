# Security Implementation Guide

## 🔒 Security Features Implemented

### 1. **Authentication & Authorization**
- **JWT Tokens**: Secure token-based authentication
- **Role-based Access Control**: Customer, Seller, Admin roles
- **Password Security**: BCrypt hashing with configurable cost
- **Session Management**: Redis-based session storage

### 2. **Password Security**
- **Strong Password Policy**: 
  - Minimum 12 characters
  - Must contain uppercase, lowercase, numbers, and special characters
  - Forbidden common passwords and patterns
  - No sequential characters (abc, 123, etc.)
- **Password Validation**: Real-time strength checking
- **Secure Storage**: BCrypt with cost factor 12

### 3. **Rate Limiting**
- **General API**: 100 requests/minute per IP
- **Authentication Endpoints**: 30 requests/minute per IP (stricter)
- **Burst Protection**: Configurable burst allowance
- **Rate Limit Headers**: X-RateLimit-* headers for client information

### 4. **Security Headers**
- **X-Content-Type-Options**: Prevents MIME sniffing
- **X-Frame-Options**: Prevents clickjacking
- **X-XSS-Protection**: XSS filtering
- **Content-Security-Policy**: Prevents code injection
- **Strict-Transport-Security**: HTTPS enforcement
- **Referrer-Policy**: Controls referrer information

### 5. **Input Validation & SQL Injection Prevention**
- **Parameterized Queries**: All database queries use GORM's parameterization
- **Input Validation**: Comprehensive struct validation
- **XSS Prevention**: Output encoding and CSP headers
- **File Upload Security**: Type validation and size limits

### 6. **Infrastructure Security**
- **Environment Variables**: Sensitive data in environment variables
- **Docker Security**: No hardcoded secrets in containers
- **Database Security**: Strong passwords and SSL support
- **HTTPS Enforcement**: Automatic HTTP to HTTPS redirect in production

## 🚀 Security Configuration

### Environment Variables
```env
# Security Settings
JWT_SECRET=your_very_long_and_secure_jwt_secret_here
BCRYPT_COST=12
RATE_LIMIT_REQUESTS=100
MAX_LOGIN_ATTEMPTS=5
LOCKOUT_DURATION=30m

# Database Security
DB_PASSWORD=your_secure_database_password
DB_SSL_MODE=require  # Enable in production

# Admin Account
ADMIN_EMAIL=admin@yourdomain.com
ADMIN_PASSWORD=YourSecureAdminPassword123!
```

### Docker Compose Security
```yaml
# Use environment variables instead of hardcoded values
environment:
  DB_PASSWORD: ${DB_PASSWORD}
  JWT_SECRET: ${JWT_SECRET}
```

### Password Policy
```go
// Password Requirements:
- Minimum 12 characters
- At least 1 uppercase letter
- At least 1 lowercase letter  
- At least 1 number
- At least 1 special character
- No common passwords (password, admin, etc.)
- No sequential characters (abc, 123, qwerty)
- No more than 2 consecutive identical characters
```

## 🛡️ Security Best Practices

### 1. **Development**
```bash
# Use strong secrets
openssl rand -base64 64  # Generate JWT secret

# Check for vulnerabilities
go mod tidy
go list -json -m all | nancy sleuth

# Test password validation
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"password": "weak"}'  # Should fail
```

### 2. **Production Deployment**
```env
# Production settings
APP_ENV=production
DEBUG=false
DB_SSL_MODE=require
HTTPS_ONLY=true
RATE_LIMIT_REQUESTS=60  # More restrictive in production
```

### 3. **Monitoring & Alerts**
- Monitor failed login attempts
- Track rate limit violations
- Alert on security header bypass attempts
- Log password reset requests

## 🔍 Security Testing

### Rate Limiting Test
```bash
# Test rate limiting
for i in {1..35}; do
  curl -X POST http://localhost:8080/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email":"test@example.com","password":"wrong"}' &
done
# Should get 429 Too Many Requests after 30 attempts
```

### Password Strength Test
```bash
# Test weak password (should fail)
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Test",
    "last_name": "User", 
    "email": "test@example.com",
    "password": "password123",
    "role": "customer"
  }'

# Test strong password (should succeed)
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Test",
    "last_name": "User",
    "email": "test@example.com", 
    "password": "MySecure!Pass123",
    "role": "customer"
  }'
```

### Security Headers Test
```bash
# Check security headers
curl -I http://localhost:8080/health
# Should include X-Content-Type-Options, X-Frame-Options, etc.
```

## ⚠️ Security Warnings

### Critical Actions Required:
1. **Change Default Passwords**: Update all default passwords in `.env`
2. **Generate JWT Secret**: Use `openssl rand -base64 64` for production
3. **Enable HTTPS**: Configure SSL certificates in production
4. **Database SSL**: Enable `DB_SSL_MODE=require` in production
5. **Admin Account**: Create admin user with strong credentials

### Regular Security Tasks:
- [ ] Rotate JWT secrets monthly
- [ ] Review and update password policies
- [ ] Monitor rate limiting effectiveness
- [ ] Update dependencies for security patches
- [ ] Review access logs for suspicious activity
- [ ] Test backup and recovery procedures

## 🔧 Security Headers Configuration

The following security headers are automatically added:

| Header | Value | Purpose |
|--------|--------|---------|
| `X-Content-Type-Options` | `nosniff` | Prevent MIME sniffing |
| `X-Frame-Options` | `DENY` | Prevent clickjacking |
| `X-XSS-Protection` | `1; mode=block` | Enable XSS filtering |
| `Content-Security-Policy` | Restrictive policy | Prevent code injection |
| `Strict-Transport-Security` | `max-age=31536000` | Enforce HTTPS |
| `Referrer-Policy` | `strict-origin-when-cross-origin` | Control referrer info |

## 📊 Security Monitoring

### Metrics to Track:
- Failed login attempts per IP
- Rate limit violations
- Password reset requests
- Admin access patterns
- API error rates
- File upload attempts

### Log Analysis:
```bash
# Monitor failed logins
grep "authentication failed" /var/log/app.log

# Check rate limit violations  
grep "rate limit exceeded" /var/log/app.log

# Monitor admin access
grep "admin" /var/log/app.log | grep "success"
```

## 🚨 Incident Response

### Suspected Breach:
1. Immediately rotate JWT secrets
2. Force logout all users
3. Review access logs
4. Check for data modifications
5. Notify affected users if needed

### Rate Limit Abuse:
1. Check source IP patterns
2. Consider IP blocking
3. Adjust rate limits if needed
4. Monitor for distributed attacks

This security implementation provides defense-in-depth protection for your e-commerce API.
