// Package email provides examples of how to use the email service
package email

import (
	"fmt"
)

// ExampleSendWelcomeEmail demonstrates how to send a welcome email
func ExampleSendWelcomeEmail() {
	// Initialize the email service
	emailService := NewEmailService()

	// Check if SMTP is configured
	if emailService.host == "" {
		fmt.Println("SMTP not configured. Set SMTP environment variables to send emails.")
		return
	}

	// Prepare email data
	to := "user@example.com"
	data := map[string]interface{}{
		"Name":           "John Doe",
		"MarketplaceURL": "https://perthmarketplace.com",
	}

	// Send welcome email
	err := emailService.SendWelcome(to, data)
	if err != nil {
		fmt.Printf("Failed to send welcome email: %v\n", err)
		return
	}

	fmt.Println("Welcome email sent successfully!")
}

// ExampleSendOrderConfirmation demonstrates how to send an order confirmation email
func ExampleSendOrderConfirmation() {
	// Initialize the email service
	emailService := NewEmailService()

	// Check if SMTP is configured
	if emailService.host == "" {
		fmt.Println("SMTP not configured. Set SMTP environment variables to send emails.")
		return
	}

	// Prepare order items
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

	// Prepare email data
	to := "customer@example.com"
	orderID := "ORD-2023-001"
	data := map[string]interface{}{
		"CustomerName":    "Jane Smith",
		"OrderID":         orderID,
		"OrderDate":       "January 15, 2023",
		"TotalAmount":     "₦30,000.00",
		"Items":           items,
		"OrderDetailsURL": "https://perthmarketplace.com/orders/ORD-2023-001",
		"MarketplaceURL":  "https://perthmarketplace.com",
	}

	// Send order confirmation email
	err := emailService.SendOrderConfirmation(to, orderID, data)
	if err != nil {
		fmt.Printf("Failed to send order confirmation email: %v\n", err)
		return
	}

	fmt.Println("Order confirmation email sent successfully!")
}

// ExampleSendMerchantOrderNotification demonstrates how to send a merchant order notification
func ExampleSendMerchantOrderNotification() {
	// Initialize the email service
	emailService := NewEmailService()

	// Check if SMTP is configured
	if emailService.host == "" {
		fmt.Println("SMTP not configured. Set SMTP environment variables to send emails.")
		return
	}

	// Prepare order items
	items := []map[string]interface{}{
		{
			"Name":     "Handcrafted Necklace",
			"Quantity": 1,
			"Price":    "₦8,500.00",
		},
	}

	// Prepare email data
	to := "merchant@example.com"
	orderID := "ORD-2023-001"
	data := map[string]interface{}{
		"MerchantName":         "Artisan Jewelry Store",
		"OrderID":              orderID,
		"OrderDate":            "January 15, 2023",
		"TotalAmount":          "₦8,500.00",
		"Items":                items,
		"MerchantDashboardURL": "https://perthmarketplace.com/merchant/dashboard",
	}

	// Send merchant order notification
	err := emailService.SendMerchantOrderNotification(to, orderID, data)
	if err != nil {
		fmt.Printf("Failed to send merchant order notification: %v\n", err)
		return
	}

	fmt.Println("Merchant order notification sent successfully!")
}

// ExampleSendPasswordReset demonstrates how to send a password reset email
func ExampleSendPasswordReset() {
	// Initialize the email service
	emailService := NewEmailService()

	// Check if SMTP is configured
	if emailService.host == "" {
		fmt.Println("SMTP not configured. Set SMTP environment variables to send emails.")
		return
	}

	// Prepare email data
	to := "user@example.com"
	resetLink := "https://perthmarketplace.com/reset-password/abc123xyz"
	data := map[string]interface{}{
		"Name":      "John Doe",
		"ResetLink": resetLink,
	}

	// Send password reset email
	err := emailService.SendPasswordReset(to, data)
	if err != nil {
		fmt.Printf("Failed to send password reset email: %v\n", err)
		return
	}

	fmt.Println("Password reset email sent successfully!")
}

// ExampleSendPayoutRequest demonstrates how to send a payout request confirmation email
func ExampleSendPayoutRequest() {
	// Initialize the email service
	emailService := NewEmailService()

	// Check if SMTP is configured
	if emailService.host == "" {
		fmt.Println("SMTP not configured. Set SMTP environment variables to send emails.")
		return
	}

	// Prepare email data
	to := "merchant@example.com"
	data := map[string]interface{}{
		"MerchantName": "Artisan Jewelry Store",
		"RequestID":    "PAYOUT-2023-001",
		"Amount":       "₦50,000.00",
		"Date":         "January 15, 2023",
		"Status":       "Processing",
	}

	// Send payout request email
	err := emailService.SendPayoutRequest(to, data)
	if err != nil {
		fmt.Printf("Failed to send payout request email: %v\n", err)
		return
	}

	fmt.Println("Payout request email sent successfully!")
}

// ExampleSendPayoutCompleted demonstrates how to send a payout completed notification email
func ExampleSendPayoutCompleted() {
	// Initialize the email service
	emailService := NewEmailService()

	// Check if SMTP is configured
	if emailService.host == "" {
		fmt.Println("SMTP not configured. Set SMTP environment variables to send emails.")
		return
	}

	// Prepare email data
	to := "merchant@example.com"
	data := map[string]interface{}{
		"MerchantName":  "Artisan Jewelry Store",
		"RequestID":     "PAYOUT-2023-001",
		"Amount":        "₦50,000.00",
		"CompletedDate": "January 16, 2023",
		"TransactionID": "TXN-2023-999999",
	}

	// Send payout completed email
	err := emailService.SendPayoutCompleted(to, data)
	if err != nil {
		fmt.Printf("Failed to send payout completed email: %v\n", err)
		return
	}

	fmt.Println("Payout completed email sent successfully!")
}

// ExampleSendDisputeOpened demonstrates how to send a dispute opened notification email
func ExampleSendDisputeOpened() {
	// Initialize the email service
	emailService := NewEmailService()

	// Check if SMTP is configured
	if emailService.host == "" {
		fmt.Println("SMTP not configured. Set SMTP environment variables to send emails.")
		return
	}

	// Prepare email data
	to := "customer@example.com"
	data := map[string]interface{}{
		"UserName":   "Jane Smith",
		"OrderID":    "ORD-2023-001",
		"DisputeID":  "DISPUTE-2023-001",
		"Date":       "January 17, 2023",
		"Reason":     "Defective product",
		"DisputeURL": "https://perthmarketplace.com/disputes/DISPUTE-2023-001",
	}

	// Send dispute opened email
	err := emailService.SendDisputeOpened(to, data)
	if err != nil {
		fmt.Printf("Failed to send dispute opened email: %v\n", err)
		return
	}

	fmt.Println("Dispute opened email sent successfully!")
}

// ExampleSendDisputeResolved demonstrates how to send a dispute resolved notification email
func ExampleSendDisputeResolved() {
	// Initialize the email service
	emailService := NewEmailService()

	// Check if SMTP is configured
	if emailService.host == "" {
		fmt.Println("SMTP not configured. Set SMTP environment variables to send emails.")
		return
	}

	// Prepare email data
	to := "customer@example.com"
	data := map[string]interface{}{
		"UserName":          "Jane Smith",
		"OrderID":           "ORD-2023-001",
		"DisputeID":         "DISPUTE-2023-001",
		"ResolvedDate":      "January 18, 2023",
		"Resolution":        "Refund Issued",
		"ResolutionDetails": "Your dispute has been resolved in your favor. A full refund of ₦8,500.00 has been issued to your original payment method.",
		"DisputeURL":        "https://perthmarketplace.com/disputes/DISPUTE-2023-001",
	}

	// Send dispute resolved email
	err := emailService.SendDisputeResolved(to, data)
	if err != nil {
		fmt.Printf("Failed to send dispute resolved email: %v\n", err)
		return
	}

	fmt.Println("Dispute resolved email sent successfully!")
}