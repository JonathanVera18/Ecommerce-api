package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/JonathanVera18/ecommerce-api/internal/config"
	"github.com/JonathanVera18/ecommerce-api/internal/models"
	"github.com/JonathanVera18/ecommerce-api/internal/repository"
	"github.com/JonathanVera18/ecommerce-api/internal/utils"
	"gorm.io/gorm"
)

type authService struct {
	userRepo   repository.UserRepository
	jwtService *utils.JWTService
	redis      *redis.Client
	config     *config.Config
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repository.UserRepository, cfg *config.Config, redisClient *redis.Client) AuthService {
	jwtService := utils.NewJWTService(cfg.JWT.Secret, cfg.JWT.Expiry)
	
	return &authService{
		userRepo:   userRepo,
		jwtService: jwtService,
		redis:      redisClient,
		config:     cfg,
	}
}

func (s *authService) Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Create new user
	user := &models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Role:      req.Role,
		Phone:     req.Phone,
		IsActive:  true,
	}

	// Validate password strength
	if err := utils.ValidatePassword(req.Password); err != nil {
		return nil, fmt.Errorf("password validation failed: %w", err)
	}

	// Hash password
	if err := user.HashPassword(req.Password); err != nil {
		return nil, err
	}

	// Save user
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generate JWT token
	token, err := s.jwtService.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	// Update last login
	s.userRepo.UpdateLastLogin(ctx, user.ID)

	return &models.AuthResponse{
		User:  user.ToResponse(),
		Token: token,
	}, nil
}

func (s *authService) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	// Verify password
	if err := user.CheckPassword(req.Password); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := s.jwtService.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	// Update last login
	s.userRepo.UpdateLastLogin(ctx, user.ID)

	return &models.AuthResponse{
		User:  user.ToResponse(),
		Token: token,
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, token string) (string, error) {
	return s.jwtService.RefreshToken(token)
}

func (s *authService) Logout(ctx context.Context, userID uint) error {
	// In a more complex implementation, you might blacklist the token in Redis
	// For now, we'll just return nil as JWT tokens are stateless
	return nil
}

func (s *authService) GetCurrentUser(ctx context.Context, userID uint) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *authService) ChangePassword(ctx context.Context, userID uint, req *models.PasswordChangeRequest) error {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify current password
	if err := user.CheckPassword(req.CurrentPassword); err != nil {
		return errors.New("current password is incorrect")
	}

	// Validate new password strength
	if err := utils.ValidatePassword(req.NewPassword); err != nil {
		return fmt.Errorf("new password validation failed: %w", err)
	}

	// Hash new password
	if err := user.HashPassword(req.NewPassword); err != nil {
		return err
	}

	// Update user
	return s.userRepo.Update(ctx, user)
}

func (s *authService) ValidateToken(token string) (uint, error) {
	claims, err := s.jwtService.ValidateToken(token)
	if err != nil {
		return 0, err
	}
	
	return claims.UserID, nil
}

// GetJWTService returns the JWT service instance
func (s *authService) GetJWTService() *utils.JWTService {
	return s.jwtService
}

// ForgotPassword initiates password reset process
func (s *authService) ForgotPassword(ctx context.Context, email string) error {
	// Check if user exists
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Don't reveal if email exists or not
			return nil
		}
		return err
	}

	// Generate reset token
	resetToken, err := utils.GenerateRandomToken(32)
	if err != nil {
		return err
	}

	// Create password reset token
	passwordResetToken := &models.PasswordResetToken{
		UserID:    user.ID,
		Token:     resetToken,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hours
	}

	if err := s.userRepo.CreatePasswordResetToken(ctx, passwordResetToken); err != nil {
		return err
	}

	// Here you would typically send an email with the reset token
	// For now, we'll just return success
	// TODO: Implement email sending
	
	return nil
}

// ResetPassword resets user password using token
func (s *authService) ResetPassword(ctx context.Context, token string, newPassword string) error {
	// Get token
	resetToken, err := s.userRepo.GetPasswordResetToken(ctx, token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("invalid or expired token")
		}
		return err
	}

	// Validate new password strength
	if err := utils.ValidatePassword(newPassword); err != nil {
		return fmt.Errorf("password validation failed: %w", err)
	}

	// Hash new password
	if err := resetToken.User.HashPassword(newPassword); err != nil {
		return err
	}

	// Update user password
	if err := s.userRepo.Update(ctx, &resetToken.User); err != nil {
		return err
	}

	// Mark token as used
	if err := s.userRepo.MarkPasswordResetTokenUsed(ctx, token); err != nil {
		return err
	}

	return nil
}

// VerifyEmail verifies user email using token
func (s *authService) VerifyEmail(ctx context.Context, token string) error {
	// Get token
	verifyToken, err := s.userRepo.GetEmailVerificationToken(ctx, token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("invalid or expired token")
		}
		return err
	}

	// Mark email as verified
	if err := s.userRepo.MarkEmailVerified(ctx, verifyToken.UserID); err != nil {
		return err
	}

	// Mark token as used
	if err := s.userRepo.MarkEmailVerificationTokenUsed(ctx, token); err != nil {
		return err
	}

	return nil
}

// ResendVerification resends email verification token
func (s *authService) ResendVerification(ctx context.Context, email string) error {
	// Check if user exists
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// Check if already verified
	if user.IsVerified {
		return errors.New("email already verified")
	}

	// Generate verification token
	verificationToken, err := utils.GenerateRandomToken(32)
	if err != nil {
		return err
	}

	// Create email verification token
	emailVerificationToken := &models.EmailVerificationToken{
		UserID:    user.ID,
		Token:     verificationToken,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hours
	}

	if err := s.userRepo.CreateEmailVerificationToken(ctx, emailVerificationToken); err != nil {
		return err
	}

	// Here you would typically send an email with the verification token
	// For now, we'll just return success
	// TODO: Implement email sending
	
	return nil
}
