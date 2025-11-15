package email

import (
	"fmt"
)

func TestEmailService() {
	fmt.Println("Testing email service...")
	
	// Create email service
	emailService := NewEmailService()
	
	// Check if SMTP is configured
	if emailService.host == "" {
		fmt.Println("SMTP not configured. Set SMTP environment variables to send emails.")
		fmt.Println("Required environment variables:")
		fmt.Println("  SMTP_HOST=your-smtp-host")
		fmt.Println("  SMTP_PORT=your-smtp-port")
		fmt.Println("  SMTP_USERNAME=your-username")
		fmt.Println("  SMTP_PASSWORD=your-password")
		fmt.Println("  SMTP_FROM=sender@example.com")
		return
	}
	
	fmt.Printf("Email service configured with host: %s, port: %d\n", emailService.host, emailService.port)
	fmt.Println("Email service is ready to send emails!")
}