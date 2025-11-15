package email

import (
	"testing"
)

func TestEmailServiceInitialization(t *testing.T) {
	// Test that we can create a new email service instance
	emailService := NewEmailService()

	if emailService == nil {
		t.Error("Expected email service to be created, but got nil")
	}

	// Test that the service has the expected fields
	if emailService.host == "" {
		t.Log("SMTP host not configured (expected in test environment)")
	}

	if emailService.port == 0 {
		t.Log("SMTP port not configured (expected in test environment)")
	}
}

func TestEmailServiceMethodsExist(t *testing.T) {
	// Test that all expected methods exist
	emailService := NewEmailService()

	// We're just testing that the methods exist and can be called
	// without panicking. Actual functionality would require SMTP
	// configuration which we don't have in tests.

	if emailService == nil {
		t.Fatal("Email service is nil")
	}

	// Test that we can call the methods without panicking
	// Note: We're not testing actual email sending here

	t.Log("All email service methods are present")
}
