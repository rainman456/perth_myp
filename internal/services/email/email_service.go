package email

import (
	"bytes"
	"crypto/tls"
	"embed"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"strconv"
)

//go:embed templates/*.html
var templateFS embed.FS

// EmailService handles sending emails with HTML templates
type EmailService struct {
	host     string
	port     int
	username string
	password string
	from     string
}

// EmailData contains data for email templates
type EmailData struct {
	Subject string
	To      string
	Data    map[string]any
}

// NewEmailService creates a new email service instance
func NewEmailService() *EmailService {
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	return &EmailService{
		host:     os.Getenv("SMTP_HOST"),
		port:     port,
		username: os.Getenv("SMTP_USERNAME"),
		password: os.Getenv("SMTP_PASSWORD"),
		from:     os.Getenv("SMTP_FROM"),
	}
}

// SendEmail sends an email with the specified template
func (e *EmailService) SendEmail(to, subject, templateName string, data map[string]any) error {
	// Read the template content from embedded FS
	tmplPath := fmt.Sprintf("templates/%s.html", templateName)
	content, err := templateFS.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("failed to read template: %w", err)
	}

	// Parse the template content
	tmpl := template.New(templateName)
	tmpl, err = tmpl.Parse(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute the template with data
	var tpl bytes.Buffer
	if err := tmpl.ExecuteTemplate(&tpl, "content", data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Create the email content (full message with headers)
	emailContent := fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=utf-8\r\n"+
			"\r\n"+
			"%s",
		e.from, to, subject, tpl.String())

	// Set up TLS config
	tlsConfig := &tls.Config{
		ServerName: e.host,
		// For testing, you can temporarily set InsecureSkipVerify: true if you hit cert issues,
		// but set it to false in production for security.
		InsecureSkipVerify: false,
	}

	// Dial the server with TLS
	addr := fmt.Sprintf("%s:%d", e.host, e.port)
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to dial TLS: %w", err)
	}
	defer conn.Close()

	// Create SMTP client
	client, err := smtp.NewClient(conn, e.host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	// Authenticate
	auth := smtp.PlainAuth("", e.username, e.password, e.host)
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	// Set from
	if err = client.Mail(e.from); err != nil {
		return fmt.Errorf("failed to set from: %w", err)
	}

	// Set recipient
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	// Send data
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to start data: %w", err)
	}
	_, err = writer.Write([]byte(emailContent))
	if err != nil {
		return fmt.Errorf("failed to write email content: %w", err)
	}
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	// Quit
	err = client.Quit()
	if err != nil {
		return fmt.Errorf("failed to quit: %w", err)
	}

	return nil
}

// SendOrderConfirmation sends an order confirmation email
func (e *EmailService) SendOrderConfirmation(to, orderID string, data map[string]interface{}) error {
	subject := fmt.Sprintf("Order Confirmation - #%s", orderID)
	return e.SendEmail(to, subject, "order_confirmation", data)
}

// SendOrderStatusUpdate sends an order status update email
func (e *EmailService) SendOrderStatusUpdate(to, orderID string, data map[string]interface{}) error {
	subject := fmt.Sprintf("Order Status Update - #%s", orderID)
	return e.SendEmail(to, subject, "order_status_update", data)
}

// SendMerchantOrderNotification sends a notification to merchant about new order
func (e *EmailService) SendMerchantOrderNotification(to, orderID string, data map[string]interface{}) error {
	subject := fmt.Sprintf("New Order Received - #%s", orderID)
	return e.SendEmail(to, subject, "merchant_order_notification", data)
}

// SendPasswordReset sends a password reset email
func (e *EmailService) SendPasswordReset(to string, data map[string]interface{}) error {
	subject := "Password Reset Request"
	return e.SendEmail(to, subject, "password_reset", data)
}

// SendWelcome sends a welcome email
func (e *EmailService) SendWelcome(to string, data map[string]interface{}) error {
	subject := "Welcome to Our Platform"
	return e.SendEmail(to, subject, "welcome", data)
}

// SendPayoutRequest sends a payout request confirmation email
func (e *EmailService) SendPayoutRequest(to string, data map[string]interface{}) error {
	subject := "Payout Request Submitted"
	return e.SendEmail(to, subject, "payout_request", data)
}

// SendPayoutCompleted sends a payout completed notification email
func (e *EmailService) SendPayoutCompleted(to string, data map[string]interface{}) error {
	subject := "Payout Completed"
	return e.SendEmail(to, subject, "payout_completed", data)
}

// SendDisputeOpened sends a dispute opened notification email
func (e *EmailService) SendDisputeOpened(to string, data map[string]interface{}) error {
	subject := "Dispute Opened for Your Order"
	return e.SendEmail(to, subject, "dispute_opened", data)
}

// SendDisputeResolved sends a dispute resolved notification email
func (e *EmailService) SendDisputeResolved(to string, data map[string]interface{}) error {
	subject := "Dispute Resolved"
	return e.SendEmail(to, subject, "dispute_resolved", data)
}