package payment

import (
	"fmt"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/JonathanVera18/ecommerce-api/internal/config"
	"github.com/JonathanVera18/ecommerce-api/internal/models"
)

type stripeService struct {
	config *config.Config
}

// NewStripeService creates a new Stripe payment service
func NewStripeService(cfg *config.Config) Service {
	stripe.Key = cfg.Stripe.SecretKey
	
	return &stripeService{
		config: cfg,
	}
}

func (s *stripeService) CreatePaymentIntent(req *models.PaymentRequest) (string, error) {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(req.Amount * 100)), // Convert to cents
		Currency: stripe.String(req.Currency),
		Metadata: map[string]string{
			"order_id": fmt.Sprintf("%d", req.OrderID),
		},
	}
	
	if req.PaymentMethodID != nil {
		params.PaymentMethod = req.PaymentMethodID
		params.ConfirmationMethod = stripe.String("manual")
		params.Confirm = stripe.Bool(true)
	}
	
	pi, err := paymentintent.New(params)
	if err != nil {
		return "", err
	}
	
	return pi.ID, nil
}

func (s *stripeService) ConfirmPayment(paymentIntentID string) error {
	params := &stripe.PaymentIntentConfirmParams{}
	_, err := paymentintent.Confirm(paymentIntentID, params)
	return err
}

func (s *stripeService) RefundPayment(paymentIntentID string, amount float64) error {
	// Implementation would depend on your refund requirements
	// This is a placeholder
	return nil
}

func (s *stripeService) GetPayment(paymentIntentID string) (*PaymentInfo, error) {
	pi, err := paymentintent.Get(paymentIntentID, nil)
	if err != nil {
		return nil, err
	}
	
	return &PaymentInfo{
		ID:       pi.ID,
		Amount:   float64(pi.Amount) / 100, // Convert from cents
		Currency: string(pi.Currency),
		Status:   string(pi.Status),
	}, nil
}
