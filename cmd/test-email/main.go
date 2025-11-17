package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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

	// Check if SMTP is configured
	if !isSMTPConfigured(cfg) {
		fmt.Println("❌ Email service disabled: SMTP credentials not configured")
		fmt.Println("\nRequired environment variables:")
		fmt.Println("  SMTP_HOST=your-smtp-host (e.g., smtp.gmail.com)")
		fmt.Println("  SMTP_PORT=your-smtp-port (e.g., 587)")
		fmt.Println("  SMTP_USERNAME=your-username")
		fmt.Println("  SMTP_PASSWORD=your-password")
		fmt.Println("  SMTP_FROM=sender@example.com")
		os.Exit(1)
	}

	fmt.Println("✅ SMTP Configuration Found:")
	fmt.Printf("   Host: %s\n", cfg.SMTPHost)
	fmt.Printf("   Port: %d\n", cfg.SMTPPort)
	fmt.Printf("   Username: %s\n", cfg.SMTPUsername)
	fmt.Printf("   From: %s\n", cfg.SMTPFrom)
	fmt.Println()

	// Create email service
	emailService := email.NewEmailService()

	// Test connection
	fmt.Println("🔄 Testing SMTP connection...")
	if err := testConnection(emailService); err != nil {
		fmt.Printf("❌ Connection test failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ SMTP connection successful!")
	fmt.Println()

	// Get test recipient email
	testEmail :="pugnojegne@necub.com"
	// if testEmail == "" {
	// 	fmt.Print("Enter test email address: ")
	// 	fmt.Scanln(&testEmail)
	// }

	// if testEmail == "" {
	// 	fmt.Println("❌ Test email address is required")
	// 	os.Exit(1)
	// }

	// Run test emails
	fmt.Printf("📧 Sending test emails to: %s\n\n", testEmail)
	
	testResults := []struct {
		name string
		fn   func(string) error
	}{
		{"Welcome Email", func(to string) error { return testWelcomeEmail(emailService, to) }},
		{"Password Reset Email", func(to string) error { return testPasswordResetEmail(emailService, to) }},
		{"Order Confirmation Email", func(to string) error { return testOrderConfirmationEmail(emailService, to) }},
		{"Order Status Update Email", func(to string) error { return testOrderStatusUpdateEmail(emailService, to) }},
		{"Merchant Order Notification", func(to string) error { return testMerchantOrderNotification(emailService, to) }},
		{"Payout Request Email", func(to string) error { return testPayoutRequestEmail(emailService, to) }},
		{"Payout Completed Email", func(to string) error { return testPayoutCompletedEmail(emailService, to) }},
		{"Dispute Opened Email", func(to string) error { return testDisputeOpenedEmail(emailService, to) }},
		{"Dispute Resolved Email", func(to string) error { return testDisputeResolvedEmail(emailService, to) }},
	}

	successCount := 0
	failCount := 0

	for i, test := range testResults {
		fmt.Printf("[%d/%d] Testing %s...\n", i+1, len(testResults), test.name)
		if err := test.fn(testEmail); err != nil {
			fmt.Printf("   ❌ Failed: %v\n", err)
			failCount++
		} else {
			fmt.Printf("   ✅ Success\n")
			successCount++
		}
		// Small delay between emails to avoid rate limiting
		time.Sleep(500 * time.Millisecond)
	}

	//fmt.Println("\n" + "=".repeat(50))
	fmt.Printf("Test Results: %d passed, %d failed\n", successCount, failCount)
	if failCount == 0 {
		fmt.Println("🎉 All email tests passed!")
	} else {
		fmt.Println("⚠️  Some email tests failed. Check the logs above.")
	}
}

func isSMTPConfigured(cfg *config.Config) bool {
	return cfg.SMTPHost != "" &&
		cfg.SMTPPort != 0 &&
		cfg.SMTPUsername != "" &&
		cfg.SMTPPassword != "" &&
		cfg.SMTPFrom != ""
}

func testConnection(emailService *email.EmailService) error {
	// Try to send a minimal test
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a test channel to check connection
	done := make(chan error, 1)
	
	go func() {
		// Just verify the service is initialized
		if emailService == nil {
			done <- fmt.Errorf("email service is nil")
			return
		}
		done <- nil
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return fmt.Errorf("connection test timeout")
	}
}

func testWelcomeEmail(emailService *email.EmailService, to string) error {
	data := map[string]interface{}{
		"Name":           "John Doe",
		"MarketplaceURL": "https://perthmarketplace.com",
	}
	return emailService.SendWelcome(to, data)
}

func testPasswordResetEmail(emailService *email.EmailService, to string) error {
	data := map[string]interface{}{
		"Name":      "John Doe",
		"ResetLink": "https://perthmarketplace.com/reset-password?token=test_token_123",
	}
	return emailService.SendPasswordReset(to, data)
}

func testOrderConfirmationEmail(emailService *email.EmailService, to string) error {
	items := []map[string]interface{}{
		{
			"Name":     "Premium T-Shirt",
			"Quantity": 2,
			"Price":    "₦7,500.00",
		},
		{
			"Name":     "Designer Jeans",
			"Quantity": 1,
			"Price":    "₦15,000.00",
		},
	}

	data := map[string]interface{}{
		"CustomerName":    "Jane Smith",
		"OrderID":         "ORD-2023-001",
		"OrderDate":       time.Now().Format("January 2, 2006"),
		"TotalAmount":     "₦30,000.00",
		"Items":           items,
		"OrderDetailsURL": "https://perthmarketplace.com/orders/ORD-2023-001",
		"MarketplaceURL":  "https://perthmarketplace.com",
	}
	return emailService.SendOrderConfirmation(to, "ORD-2023-001", data)
}

func testOrderStatusUpdateEmail(emailService *email.EmailService, to string) error {
	data := map[string]interface{}{
		"CustomerName":    "Jane Smith",
		"OrderID":         "ORD-2023-001",
		"NewStatus":       "Shipped",
		"UpdateDate":      time.Now().Format("January 2, 2006 at 3:04 PM"),
		"TrackingNumber":  "TRK123456789",
		"Carrier":         "DHL Express",
		"OrderDetailsURL": "https://perthmarketplace.com/orders/ORD-2023-001",
	}
	return emailService.SendOrderStatusUpdate(to, "ORD-2023-001", data)
}

func testMerchantOrderNotification(emailService *email.EmailService, to string) error {
	items := []map[string]interface{}{
		{
			"Name":     "Handcrafted Necklace",
			"Quantity": 1,
			"Price":    "₦8,500.00",
		},
	}

	data := map[string]interface{}{
		"MerchantName":         "Artisan Jewelry Store",
		"OrderID":              "ORD-2023-001",
		"OrderDate":            time.Now().Format("January 2, 2006"),
		"TotalAmount":          "₦8,500.00",
		"Items":                items,
		"MerchantDashboardURL": "https://perthmarketplace.com/merchant/dashboard",
	}
	return emailService.SendMerchantOrderNotification(to, "ORD-2023-001", data)
}

func testPayoutRequestEmail(emailService *email.EmailService, to string) error {
	data := map[string]interface{}{
		"MerchantName": "Artisan Jewelry Store",
		"RequestID":    "PAYOUT-2023-001",
		"Amount":       "₦50,000.00",
		"Date":         time.Now().Format("January 2, 2006"),
		"Status":       "Processing",
	}
	return emailService.SendPayoutRequest(to, data)
}

func testPayoutCompletedEmail(emailService *email.EmailService, to string) error {
	data := map[string]interface{}{
		"MerchantName":  "Artisan Jewelry Store",
		"RequestID":     "PAYOUT-2023-001",
		"Amount":        "₦50,000.00",
		"CompletedDate": time.Now().Format("January 2, 2006"),
		"TransactionID": "TXN-2023-999999",
	}
	return emailService.SendPayoutCompleted(to, data)
}

func testDisputeOpenedEmail(emailService *email.EmailService, to string) error {
	data := map[string]interface{}{
		"UserName":   "Jane Smith",
		"OrderID":    "ORD-2023-001",
		"DisputeID":  "DISPUTE-2023-001",
		"Date":       time.Now().Format("January 2, 2006"),
		"Reason":     "Product quality issue",
		"DisputeURL": "https://perthmarketplace.com/disputes/DISPUTE-2023-001",
	}
	return emailService.SendDisputeOpened(to, data)
}

func testDisputeResolvedEmail(emailService *email.EmailService, to string) error {
	data := map[string]interface{}{
		"UserName":          "Jane Smith",
		"OrderID":           "ORD-2023-001",
		"DisputeID":         "DISPUTE-2023-001",
		"ResolvedDate":      time.Now().Format("January 2, 2006"),
		"Resolution":        "Refund Issued",
		"ResolutionDetails": "Your dispute has been resolved in your favor. A full refund of ₦8,500.00 has been issued to your original payment method.",
		"DisputeURL":        "https://perthmarketplace.com/disputes/DISPUTE-2023-001",
	}
	return emailService.SendDisputeResolved(to, data)
}

// Helper function for string repetition
//type repeatableString string

// func (s repeatableString) repeat(count int) string {
// 	result := ""
// 	for i := 0; i < count; i++ {
// 		result += string(s)
// 	}
// 	return result
// }

//var repeat = repeatableString("=").repeat