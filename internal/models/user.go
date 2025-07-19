package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// UserRole represents user roles
type UserRole string

const (
	RoleCustomer UserRole = "customer"
	RoleSeller   UserRole = "seller"
	RoleAdmin    UserRole = "admin"
)

// User represents a user in the system
type User struct {
	BaseModel
	FirstName    string    `json:"first_name" gorm:"type:varchar(100);not null" validate:"required,min=2,max=100"`
	LastName     string    `json:"last_name" gorm:"type:varchar(100);not null" validate:"required,min=2,max=100"`
	Email        string    `json:"email" gorm:"type:varchar(255);unique;not null" validate:"required,email"`
	Password     string    `json:"-" gorm:"type:varchar(255);not null" validate:"required,min=8"`
	Phone        *string   `json:"phone,omitempty" gorm:"type:varchar(20)" validate:"omitempty,e164"`
	Role         UserRole  `json:"role" gorm:"type:varchar(20);not null;default:'customer'" validate:"required,oneof=customer seller admin"`
	IsActive     bool      `json:"is_active" gorm:"default:true"`
	IsVerified   bool      `json:"is_verified" gorm:"default:false"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	
	// Profile information
	DateOfBirth *time.Time `json:"date_of_birth,omitempty" gorm:"type:date"`
	Gender      *string    `json:"gender,omitempty" gorm:"type:varchar(10)" validate:"omitempty,oneof=male female other"`
	Avatar      *string    `json:"avatar,omitempty" gorm:"type:varchar(500)"`
	
	// Address information
	Street     *string `json:"street,omitempty" gorm:"type:varchar(255)"`
	City       *string `json:"city,omitempty" gorm:"type:varchar(100)"`
	State      *string `json:"state,omitempty" gorm:"type:varchar(100)"`
	Country    *string `json:"country,omitempty" gorm:"type:varchar(100)"`
	PostalCode *string `json:"postal_code,omitempty" gorm:"type:varchar(20)"`
	
	// Seller specific fields
	StoreName        *string `json:"store_name,omitempty" gorm:"type:varchar(255)"`
	StoreDescription *string `json:"store_description,omitempty" gorm:"type:text"`
	TaxID           *string `json:"tax_id,omitempty" gorm:"type:varchar(50)"`
	
	// Relationships
	Products []Product `json:"products,omitempty" gorm:"foreignKey:SellerID"`
	Orders   []Order   `json:"orders,omitempty" gorm:"foreignKey:CustomerID"`
	Reviews  []Review  `json:"reviews,omitempty" gorm:"foreignKey:UserID"`
}

// UserCreateRequest represents the request to create a user
type UserCreateRequest struct {
	FirstName string   `json:"first_name" validate:"required,min=2,max=100"`
	LastName  string   `json:"last_name" validate:"required,min=2,max=100"`
	Email     string   `json:"email" validate:"required,email"`
	Password  string   `json:"password" validate:"required,min=8"`
	Phone     *string  `json:"phone,omitempty" validate:"omitempty,e164"`
	Role      UserRole `json:"role" validate:"required,oneof=customer seller admin"`
}

// UserUpdateRequest represents the request to update a user
type UserUpdateRequest struct {
	FirstName   *string    `json:"first_name,omitempty" validate:"omitempty,min=2,max=100"`
	LastName    *string    `json:"last_name,omitempty" validate:"omitempty,min=2,max=100"`
	Phone       *string    `json:"phone,omitempty" validate:"omitempty,e164"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	Gender      *string    `json:"gender,omitempty" validate:"omitempty,oneof=male female other"`
	
	// Address information
	Street     *string `json:"street,omitempty"`
	City       *string `json:"city,omitempty"`
	State      *string `json:"state,omitempty"`
	Country    *string `json:"country,omitempty"`
	PostalCode *string `json:"postal_code,omitempty"`
	
	// Seller information
	StoreName        *string `json:"store_name,omitempty"`
	StoreDescription *string `json:"store_description,omitempty"`
	TaxID           *string `json:"tax_id,omitempty"`
}

// UserResponse represents the user response (without sensitive data)
type UserResponse struct {
	ID          uint       `json:"id"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Email       string     `json:"email"`
	Phone       *string    `json:"phone,omitempty"`
	Role        UserRole   `json:"role"`
	IsActive    bool       `json:"is_active"`
	IsVerified  bool       `json:"is_verified"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	
	// Profile information
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	Gender      *string    `json:"gender,omitempty"`
	Avatar      *string    `json:"avatar,omitempty"`
	
	// Address information
	Street     *string `json:"street,omitempty"`
	City       *string `json:"city,omitempty"`
	State      *string `json:"state,omitempty"`
	Country    *string `json:"country,omitempty"`
	PostalCode *string `json:"postal_code,omitempty"`
	
	// Seller information
	StoreName        *string `json:"store_name,omitempty"`
	StoreDescription *string `json:"store_description,omitempty"`
}

// LoginRequest represents the login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RegisterRequest represents the registration request
type RegisterRequest struct {
	FirstName string   `json:"first_name" validate:"required,min=2,max=100"`
	LastName  string   `json:"last_name" validate:"required,min=2,max=100"`
	Email     string   `json:"email" validate:"required,email"`
	Password  string   `json:"password" validate:"required,min=8"`
	Phone     *string  `json:"phone,omitempty" validate:"omitempty,e164"`
	Role      UserRole `json:"role" validate:"required,oneof=customer seller"`
}

// PasswordChangeRequest represents the password change request
type PasswordChangeRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

// HashPassword hashes a plain text password
func (u *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword verifies a password against the hash
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:               u.ID,
		FirstName:        u.FirstName,
		LastName:         u.LastName,
		Email:            u.Email,
		Phone:            u.Phone,
		Role:             u.Role,
		IsActive:         u.IsActive,
		IsVerified:       u.IsVerified,
		LastLoginAt:      u.LastLoginAt,
		CreatedAt:        u.CreatedAt,
		UpdatedAt:        u.UpdatedAt,
		DateOfBirth:      u.DateOfBirth,
		Gender:           u.Gender,
		Avatar:           u.Avatar,
		Street:           u.Street,
		City:             u.City,
		State:            u.State,
		Country:          u.Country,
		PostalCode:       u.PostalCode,
		StoreName:        u.StoreName,
		StoreDescription: u.StoreDescription,
	}
}

// FullName returns the full name of the user
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// IsCustomer checks if the user is a customer
func (u *User) IsCustomer() bool {
	return u.Role == RoleCustomer
}

// IsSeller checks if the user is a seller
func (u *User) IsSeller() bool {
	return u.Role == RoleSeller
}

// IsAdmin checks if the user is an admin
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// UserStatsResponse represents user statistics
type UserStatsResponse struct {
	TotalUsers     int64 `json:"total_users"`
	ActiveUsers    int64 `json:"active_users"`
	VerifiedUsers  int64 `json:"verified_users"`
	Customers      int64 `json:"customers"`
	Sellers        int64 `json:"sellers"`
	Admins         int64 `json:"admins"`
	NewUsersToday  int64 `json:"new_users_today"`
	NewUsersWeek   int64 `json:"new_users_week"`
	NewUsersMonth  int64 `json:"new_users_month"`
}

// ForgotPasswordRequest represents the forgot password request
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordRequest represents the reset password request
type ResetPasswordRequest struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

// VerifyEmailRequest represents the verify email request
type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

// ResendVerificationRequest represents the resend verification request
type ResendVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// PasswordResetToken represents a password reset token
type PasswordResetToken struct {
	BaseModel
	UserID    uint       `json:"user_id" gorm:"not null;index"`
	Token     string     `json:"token" gorm:"type:varchar(255);not null;unique"`
	ExpiresAt time.Time  `json:"expires_at" gorm:"not null"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
	
	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// EmailVerificationToken represents an email verification token
type EmailVerificationToken struct {
	BaseModel
	UserID    uint       `json:"user_id" gorm:"not null;index"`
	Token     string     `json:"token" gorm:"type:varchar(255);not null;unique"`
	ExpiresAt time.Time  `json:"expires_at" gorm:"not null"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
	
	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
