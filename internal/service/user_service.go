package service

import (
	"context"
	"errors"

	"github.com/JonathanVera18/ecommerce-api/internal/models"
	"github.com/JonathanVera18/ecommerce-api/internal/repository"
	"gorm.io/gorm"
)

type userService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetProfile(ctx context.Context, userID uint) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID uint, req *models.UserUpdateRequest) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}
	if req.DateOfBirth != nil {
		user.DateOfBirth = req.DateOfBirth
	}
	if req.Gender != nil {
		user.Gender = req.Gender
	}
	if req.Street != nil {
		user.Street = req.Street
	}
	if req.City != nil {
		user.City = req.City
	}
	if req.State != nil {
		user.State = req.State
	}
	if req.Country != nil {
		user.Country = req.Country
	}
	if req.PostalCode != nil {
		user.PostalCode = req.PostalCode
	}
	if req.StoreName != nil && user.IsSeller() {
		user.StoreName = req.StoreName
	}
	if req.StoreDescription != nil && user.IsSeller() {
		user.StoreDescription = req.StoreDescription
	}
	if req.TaxID != nil && user.IsSeller() {
		user.TaxID = req.TaxID
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *userService) GetUsers(ctx context.Context, page, limit int, role *models.UserRole) ([]models.UserResponse, int64, error) {
	users, total, err := s.userRepo.List(ctx, page, limit, role)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]models.UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	return responses, total, nil
}

func (s *userService) GetUserByID(ctx context.Context, id uint) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *userService) CreateUser(ctx context.Context, req *models.UserCreateRequest) (*models.UserResponse, error) {
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

	// Hash password
	if err := user.HashPassword(req.Password); err != nil {
		return nil, err
	}

	// Save user
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *userService) UpdateUser(ctx context.Context, id uint, req *models.UserUpdateRequest) (*models.UserResponse, error) {
	return s.UpdateProfile(ctx, id, req)
}

func (s *userService) DeleteUser(ctx context.Context, id uint) error {
	return s.userRepo.Delete(ctx, id)
}

func (s *userService) GetUserStats(ctx context.Context) (*models.UserStatsResponse, error) {
	stats, err := s.userRepo.GetStats(ctx)
	if err != nil {
		return nil, err
	}
	
	// Convert to models.UserStatsResponse
	return &models.UserStatsResponse{
		TotalUsers:    stats.TotalUsers,
		ActiveUsers:   stats.ActiveUsers,
		VerifiedUsers: stats.VerifiedUsers,
		Customers:     stats.Customers,
		Sellers:       stats.Sellers,
		Admins:        stats.Admins,
		NewUsersToday: stats.NewUsersToday,
		NewUsersWeek:  stats.NewUsersWeek,
		NewUsersMonth: stats.NewUsersMonth,
	}, nil
}
