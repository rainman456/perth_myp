package main

import (
	"fmt"
	"log"
	"os"

	"api-customer-merchant/internal/config"
	"api-customer-merchant/internal/services/email"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Check SMTP configuration
	emailEnabled := cfg.SMTPHost != "" && cfg.SMTPPort != 0 && 
		cfg.SMTPUsername != "" && cfg.SMTPPassword != ""

	if !emailEnabled {
		fmt.Println("❌ Email service disabled: SMTP credentials not configured")
		fmt.Println("\nRequired environment variables:")
		fmt.Println("  SMTP_HOST=your-smtp-host")
		fmt.Println("  SMTP_PORT=your-smtp-port")
		fmt.Println("  SMTP_USERNAME=your-username")
		fmt.Println("  SMTP_PASSWORD=your-password")
		fmt.Println("  SMTP_FROM=sender@example.com")
		return
	}

	fmt.Println("✅ Email Configuration Loaded:")
	fmt.Printf("   Host: %s\n", cfg.SMTPHost)
	fmt.Printf("   Port: %d\n", cfg.SMTPPort)
	fmt.Printf("   From: %s\n", cfg.SMTPFrom)
	fmt.Println()

	// Create email service
	emailService := email.NewEmailService()
	if emailService == nil {
		fmt.Println("❌ Failed to create email service")
		return
	}

	fmt.Println("✅ Email service created successfully")
	fmt.Println()

	// Get test email
	testEmail := os.Getenv("TEST_EMAIL")
	if testEmail == "" {
		testEmail = "simonlevi453@gmail.com" // Default test email
	}

	fmt.Printf("📧 Sending test email to: %s\n", testEmail)

	// Prepare test data
	data := map[string]interface{}{
		"Name":           "Test User",
		"MarketplaceURL": "https://perthmarketplace.com",
	}

	// Send test email
	err := emailService.SendWelcome(testEmail, data)
	if err != nil {
		fmt.Printf("❌ Test failed: %v\n", err)
		return
	}

	fmt.Println("✅ Test email sent successfully!")
	fmt.Println("   Check your inbox (and spam folder) for the welcome email")
}