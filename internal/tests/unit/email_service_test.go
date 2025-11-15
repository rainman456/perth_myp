package unit

import (
	"api-customer-merchant/internal/services/email"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmailService_SendOrderConfirmation(t *testing.T) {
	// Create a new email service instance
	emailService := email.NewEmailService()

	// Prepare test data
	to := "test@example.com"
	orderID := "12345"
	data := map[string]interface{}{
		"CustomerName": "John Doe",
		"OrderID":      "12345",
		"OrderDate":    "January 1, 2023",
		"TotalAmount":  "₦10,000.00",
		"Items": []map[string]interface{}{
			{
				"Name":     "Test Product",
				"Quantity": 2,
				"Price":    "₦5,000.00",
			},
		},
		"OrderDetailsURL": "https://example.com/orders/12345",
		"MarketplaceURL":  "https://example.com",
	}

	// Test sending order confirmation email
	err := emailService.SendOrderConfirmation(to, orderID, data)

	// Note: Since we don't have actual SMTP credentials in test,
	// we're just testing that the method doesn't panic and returns an error
	// that indicates a connection issue rather than a template issue
	if err != nil {
		// We expect an error related to SMTP connection in test environment
		// but not a template parsing error
		assert.NotContains(t, err.Error(), "failed to parse template")
	}
}

func TestEmailService_SendWelcome(t *testing.T) {
	// Create a new email service instance
	emailService := email.NewEmailService()

	// Prepare test data
	to := "test@example.com"
	data := map[string]interface{}{
		"Name":           "John Doe",
		"MarketplaceURL": "https://example.com",
	}

	// Test sending welcome email
	err := emailService.SendWelcome(to, data)

	// Note: Since we don't have actual SMTP credentials in test,
	// we're just testing that the method doesn't panic and returns an error
	// that indicates a connection issue rather than a template issue
	if err != nil {
		// We expect an error related to SMTP connection in test environment
		// but not a template parsing error
		assert.NotContains(t, err.Error(), "failed to parse template")
	}
}

func TestEmailService_SendPasswordReset(t *testing.T) {
	// Create a new email service instance
	emailService := email.NewEmailService()

	// Prepare test data
	to := "test@example.com"
	data := map[string]interface{}{
		"Name":      "John Doe",
		"ResetLink": "https://example.com/reset-password/token123",
	}

	// Test sending password reset email
	err := emailService.SendPasswordReset(to, data)

	// Note: Since we don't have actual SMTP credentials in test,
	// we're just testing that the method doesn't panic and returns an error
	// that indicates a connection issue rather than a template issue
	if err != nil {
		// We expect an error related to SMTP connection in test environment
		// but not a template parsing error
		assert.NotContains(t, err.Error(), "failed to parse template")
	}
}
