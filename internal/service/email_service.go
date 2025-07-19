package service

import (
	"context"
	"fmt"

	"github.com/JonathanVera18/ecommerce-api/internal/models"
	"github.com/JonathanVera18/ecommerce-api/pkg/email"
)

type emailService struct {
	emailSender email.Service
}

func NewEmailService(emailSender email.Service) EmailService {
	return &emailService{
		emailSender: emailSender,
	}
}

func (s *emailService) SendWelcomeEmail(ctx context.Context, user *models.User) error {
	return s.emailSender.SendWelcomeEmail(user.Email, user.FirstName)
}

func (s *emailService) SendOrderConfirmationEmail(ctx context.Context, user *models.User, order *models.Order) error {
	return s.emailSender.SendOrderConfirmationEmail(user.Email, order)
}

func (s *emailService) SendOrderStatusUpdateEmail(ctx context.Context, user *models.User, order *models.Order) error {
	switch order.Status {
	case models.OrderStatusShipped:
		return s.emailSender.SendOrderShippedEmail(user.Email, order)
	case models.OrderStatusDelivered:
		return s.emailSender.SendOrderDeliveredEmail(user.Email, order)
	default:
		// For other status updates, we'll use the shipped email template for now
		return s.emailSender.SendOrderShippedEmail(user.Email, order)
	}
}

func (s *emailService) SendPasswordResetEmail(ctx context.Context, user *models.User, resetToken string) error {
	resetLink := fmt.Sprintf("https://yourdomain.com/reset-password?token=%s", resetToken)
	return s.emailSender.SendPasswordResetEmail(user.Email, resetLink)
}

func (s *emailService) SendEmailVerificationEmail(ctx context.Context, user *models.User, verificationToken string) error {
	// For now, use password reset email as template since verification email is not in the interface
	verificationLink := fmt.Sprintf("https://yourdomain.com/verify-email?token=%s", verificationToken)
	return s.emailSender.SendPasswordResetEmail(user.Email, verificationLink)
}

func (s *emailService) SendLowStockAlert(ctx context.Context, seller *models.User, product *models.Product) error {
	// Since this is not in the email.Service interface, we'll use a basic welcome email format
	return s.emailSender.SendWelcomeEmail(seller.Email, seller.FirstName)
}

func (s *emailService) SendNewReviewNotification(ctx context.Context, seller *models.User, product *models.Product, review *models.Review) error {
	// Since this is not in the email.Service interface, we'll use a basic welcome email format
	return s.emailSender.SendWelcomeEmail(seller.Email, seller.FirstName)
}
