package repository

import (
	"context"
	"time"

	"github.com/JonathanVera18/ecommerce-api/internal/models"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

func (r *userRepository) List(ctx context.Context, page, limit int, role *models.UserRole) ([]models.User, int64, error) {
	var users []models.User
	var total int64
	
	query := r.db.WithContext(ctx).Model(&models.User{})
	
	if role != nil {
		query = query.Where("role = ?", *role)
	}
	
	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Get paginated results
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}
	
	return users, total, nil
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, id uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Update("last_login_at", now).Error
}

func (r *userRepository) GetStats(ctx context.Context) (*models.UserStatsResponse, error) {
	var stats models.UserStatsResponse
	
	// Total users
	if err := r.db.WithContext(ctx).Model(&models.User{}).Count(&stats.TotalUsers).Error; err != nil {
		return nil, err
	}
	
	// Active users
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("is_active = ?", true).Count(&stats.ActiveUsers).Error; err != nil {
		return nil, err
	}
	
	// Verified users
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("is_verified = ?", true).Count(&stats.VerifiedUsers).Error; err != nil {
		return nil, err
	}
	
	// Users by role
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("role = ?", models.RoleCustomer).Count(&stats.Customers).Error; err != nil {
		return nil, err
	}
	
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("role = ?", models.RoleSeller).Count(&stats.Sellers).Error; err != nil {
		return nil, err
	}
	
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("role = ?", models.RoleAdmin).Count(&stats.Admins).Error; err != nil {
		return nil, err
	}
	
	// New users by time period
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekAgo := today.AddDate(0, 0, -7)
	monthAgo := today.AddDate(0, -1, 0)
	
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("created_at >= ?", today).Count(&stats.NewUsersToday).Error; err != nil {
		return nil, err
	}
	
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("created_at >= ?", weekAgo).Count(&stats.NewUsersWeek).Error; err != nil {
		return nil, err
	}
	
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("created_at >= ?", monthAgo).Count(&stats.NewUsersMonth).Error; err != nil {
		return nil, err
	}
	
	return &stats, nil
}

func (r *userRepository) CreatePasswordResetToken(ctx context.Context, token *models.PasswordResetToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *userRepository) GetPasswordResetToken(ctx context.Context, tokenStr string) (*models.PasswordResetToken, error) {
	var token models.PasswordResetToken
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("token = ? AND expires_at > NOW() AND used_at IS NULL", tokenStr).
		First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *userRepository) MarkPasswordResetTokenUsed(ctx context.Context, tokenStr string) error {
	return r.db.WithContext(ctx).
		Model(&models.PasswordResetToken{}).
		Where("token = ?", tokenStr).
		Update("used_at", gorm.Expr("NOW()")).Error
}

func (r *userRepository) CreateEmailVerificationToken(ctx context.Context, token *models.EmailVerificationToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *userRepository) GetEmailVerificationToken(ctx context.Context, tokenStr string) (*models.EmailVerificationToken, error) {
	var token models.EmailVerificationToken
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("token = ? AND expires_at > NOW() AND used_at IS NULL", tokenStr).
		First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *userRepository) MarkEmailVerificationTokenUsed(ctx context.Context, tokenStr string) error {
	return r.db.WithContext(ctx).
		Model(&models.EmailVerificationToken{}).
		Where("token = ?", tokenStr).
		Update("used_at", gorm.Expr("NOW()")).Error
}

func (r *userRepository) MarkEmailVerified(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Update("is_verified", true).Error
}
