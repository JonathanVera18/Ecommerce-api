package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/JonathanVera18/ecommerce-api/internal/config"
	"github.com/JonathanVera18/ecommerce-api/internal/models"
)

type smtpService struct {
	config *config.Config
	auth   smtp.Auth
}

// NewSMTPService creates a new SMTP email service
func NewSMTPService(cfg *config.Config) Service {
	auth := smtp.PlainAuth("", cfg.Email.SMTPUsername, cfg.Email.SMTPPassword, cfg.Email.SMTPHost)
	
	return &smtpService{
		config: cfg,
		auth:   auth,
	}
}

func (s *smtpService) sendEmail(to, subject, body string, isHTML bool) error {
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n")
	
	if isHTML {
		msg = append(msg, []byte("Content-Type: text/html; charset=UTF-8\r\n\r\n")...)
	} else {
		msg = append(msg, []byte("Content-Type: text/plain; charset=UTF-8\r\n\r\n")...)
	}
	
	msg = append(msg, []byte(body)...)

	addr := fmt.Sprintf("%s:%d", s.config.Email.SMTPHost, s.config.Email.SMTPPort)
	return smtp.SendMail(addr, s.auth, s.config.Email.FromEmail, []string{to}, msg)
}

func (s *smtpService) SendWelcomeEmail(to, name string) error {
	subject := "Welcome to Our E-commerce Platform!"
	body := fmt.Sprintf(`
		<html>
		<body>
			<h1>Welcome %s!</h1>
			<p>Thank you for joining our e-commerce platform. We're excited to have you as part of our community.</p>
			<p>You can now start browsing products and making purchases.</p>
			<p>If you have any questions, feel free to contact our support team.</p>
			<br>
			<p>Best regards,<br>The E-commerce Team</p>
		</body>
		</html>
	`, name)
	
	return s.sendEmail(to, subject, body, true)
}

func (s *smtpService) SendOrderConfirmationEmail(to string, order *models.Order) error {
	subject := fmt.Sprintf("Order Confirmation - Order #%s", order.OrderNumber)
	
	tmpl := `
		<html>
		<body>
			<h1>Order Confirmation</h1>
			<p>Dear {{.Customer.FirstName}},</p>
			<p>Thank you for your order! We've received your order and are processing it.</p>
			
			<h2>Order Details</h2>
			<p><strong>Order Number:</strong> {{.OrderNumber}}</p>
			<p><strong>Order Date:</strong> {{.CreatedAt.Format "January 2, 2006"}}</p>
			<p><strong>Total Amount:</strong> ${{printf "%.2f" .TotalAmount}}</p>
			
			<h3>Items Ordered</h3>
			<table border="1" style="border-collapse: collapse; width: 100%;">
				<tr>
					<th>Product</th>
					<th>Quantity</th>
					<th>Price</th>
					<th>Total</th>
				</tr>
				{{range .OrderItems}}
				<tr>
					<td>{{.ProductName}}</td>
					<td>{{.Quantity}}</td>
					<td>${{printf "%.2f" .UnitPrice}}</td>
					<td>${{printf "%.2f" .TotalPrice}}</td>
				</tr>
				{{end}}
			</table>
			
			<h3>Shipping Address</h3>
			<p>{{.ShippingFirstName}} {{.ShippingLastName}}<br>
			{{.ShippingStreet}}<br>
			{{.ShippingCity}}, {{.ShippingState}} {{.ShippingPostalCode}}<br>
			{{.ShippingCountry}}</p>
			
			<p>We'll send you another email when your order ships.</p>
			
			<p>Best regards,<br>The E-commerce Team</p>
		</body>
		</html>
	`
	
	t, err := template.New("order").Parse(tmpl)
	if err != nil {
		return err
	}
	
	var body bytes.Buffer
	if err := t.Execute(&body, order); err != nil {
		return err
	}
	
	return s.sendEmail(to, subject, body.String(), true)
}

func (s *smtpService) SendOrderShippedEmail(to string, order *models.Order) error {
	subject := fmt.Sprintf("Your Order Has Shipped - Order #%s", order.OrderNumber)
	
	trackingInfo := ""
	if order.TrackingNumber != nil {
		trackingInfo = fmt.Sprintf("<p><strong>Tracking Number:</strong> %s</p>", *order.TrackingNumber)
	}
	
	body := fmt.Sprintf(`
		<html>
		<body>
			<h1>Your Order Has Shipped!</h1>
			<p>Dear %s,</p>
			<p>Great news! Your order #%s has been shipped and is on its way to you.</p>
			
			%s
			
			<p>You can expect to receive your order within 3-7 business days.</p>
			
			<p>Thank you for your business!</p>
			
			<p>Best regards,<br>The E-commerce Team</p>
		</body>
		</html>
	`, order.ShippingFirstName, order.OrderNumber, trackingInfo)
	
	return s.sendEmail(to, subject, body, true)
}

func (s *smtpService) SendOrderDeliveredEmail(to string, order *models.Order) error {
	subject := fmt.Sprintf("Order Delivered - Order #%s", order.OrderNumber)
	body := fmt.Sprintf(`
		<html>
		<body>
			<h1>Order Delivered!</h1>
			<p>Dear %s,</p>
			<p>Your order #%s has been successfully delivered!</p>
			
			<p>We hope you love your purchase. If you have any issues with your order, please don't hesitate to contact us.</p>
			
			<p>We'd love to hear about your experience. Consider leaving a review for the products you purchased.</p>
			
			<p>Thank you for choosing us!</p>
			
			<p>Best regards,<br>The E-commerce Team</p>
		</body>
		</html>
	`, order.ShippingFirstName, order.OrderNumber)
	
	return s.sendEmail(to, subject, body, true)
}

func (s *smtpService) SendPasswordResetEmail(to, resetLink string) error {
	subject := "Password Reset Request"
	body := fmt.Sprintf(`
		<html>
		<body>
			<h1>Password Reset Request</h1>
			<p>We received a request to reset your password.</p>
			<p>Click the link below to reset your password:</p>
			<p><a href="%s">Reset Password</a></p>
			<p>If you didn't request this password reset, you can safely ignore this email.</p>
			<p>This link will expire in 1 hour.</p>
			
			<p>Best regards,<br>The E-commerce Team</p>
		</body>
		</html>
	`, resetLink)
	
	return s.sendEmail(to, subject, body, true)
}

func (s *smtpService) SendInvoiceEmail(to string, order *models.Order) error {
	subject := fmt.Sprintf("Invoice - Order #%s", order.OrderNumber)
	
	tmpl := `
		<html>
		<body>
			<h1>Invoice</h1>
			<p><strong>Order Number:</strong> {{.OrderNumber}}</p>
			<p><strong>Date:</strong> {{.CreatedAt.Format "January 2, 2006"}}</p>
			
			<h2>Bill To:</h2>
			<p>{{.ShippingFirstName}} {{.ShippingLastName}}<br>
			{{.ShippingStreet}}<br>
			{{.ShippingCity}}, {{.ShippingState}} {{.ShippingPostalCode}}<br>
			{{.ShippingCountry}}</p>
			
			<h2>Items</h2>
			<table border="1" style="border-collapse: collapse; width: 100%;">
				<tr>
					<th>Description</th>
					<th>Quantity</th>
					<th>Unit Price</th>
					<th>Total</th>
				</tr>
				{{range .OrderItems}}
				<tr>
					<td>{{.ProductName}}</td>
					<td>{{.Quantity}}</td>
					<td>${{printf "%.2f" .UnitPrice}}</td>
					<td>${{printf "%.2f" .TotalPrice}}</td>
				</tr>
				{{end}}
			</table>
			
			<h3>Summary</h3>
			<p><strong>Subtotal:</strong> ${{printf "%.2f" .SubtotalAmount}}</p>
			<p><strong>Tax:</strong> ${{printf "%.2f" .TaxAmount}}</p>
			<p><strong>Shipping:</strong> ${{printf "%.2f" .ShippingAmount}}</p>
			<p><strong>Total:</strong> ${{printf "%.2f" .TotalAmount}}</p>
			
			<p>Thank you for your business!</p>
		</body>
		</html>
	`
	
	t, err := template.New("invoice").Parse(tmpl)
	if err != nil {
		return err
	}
	
	var body bytes.Buffer
	if err := t.Execute(&body, order); err != nil {
		return err
	}
	
	return s.sendEmail(to, subject, body.String(), true)
}
