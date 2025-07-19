package email

import (
	"github.com/JonathanVera18/ecommerce-api/internal/models"
)

// Service defines the email service interface
type Service interface {
	SendWelcomeEmail(to, name string) error
	SendOrderConfirmationEmail(to string, order *models.Order) error
	SendOrderShippedEmail(to string, order *models.Order) error
	SendOrderDeliveredEmail(to string, order *models.Order) error
	SendPasswordResetEmail(to, resetLink string) error
	SendInvoiceEmail(to string, order *models.Order) error
}

// EmailTemplate represents an email template
type EmailTemplate struct {
	Subject string
	Body    string
	IsHTML  bool
}
