package payment

import "github.com/JonathanVera18/ecommerce-api/internal/models"

// Service defines the payment service interface
type Service interface {
	CreatePaymentIntent(req *models.PaymentRequest) (string, error)
	ConfirmPayment(paymentIntentID string) error
	RefundPayment(paymentIntentID string, amount float64) error
	GetPayment(paymentIntentID string) (*PaymentInfo, error)
}

// PaymentInfo represents payment information
type PaymentInfo struct {
	ID       string  `json:"id"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Status   string  `json:"status"`
}

// PaymentResult represents the result of a payment operation
type PaymentResult struct {
	Success       bool   `json:"success"`
	PaymentID     string `json:"payment_id"`
	ClientSecret  string `json:"client_secret,omitempty"`
	Error         string `json:"error,omitempty"`
}
