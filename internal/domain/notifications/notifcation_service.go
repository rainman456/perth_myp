package notifications

import (
    "context"
    "fmt"
    // Assume SMTP or Twilio lib
)

type NotificationService struct {
    // Email/SMS clients
}

func NewNotificationService() *NotificationService {
    return &NotificationService{}
}

func (s *NotificationService) SendOrderConfirmation(ctx context.Context, orderID uint, email string) error {
    // Use template: "Your order {orderID} is confirmed"
    fmt.Println("Sending email to", email) // Replace with real send
    return nil
}

func (s *NotificationService) SendMerchantNewOrder(merchantID uint, orderID uint) error {
    // Fetch merchant email, send
    return nil
}

func (s *NotificationService) SendStockAlert(merchantID uint, productID uint) error {
    // On low stock
    return nil
}