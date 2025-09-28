package payment
/*
import (
    "context"
    "crypto/hmac"
    "crypto/sha512"
    "encoding/hex"
    "encoding/json"
    "errors"
    "fmt"
    "io/ioutil"
    "net/http"
    "time"

    "api-customer-merchant/internal/config"
    "api-customer-merchant/internal/db"
    "api-customer-merchant/internal/db/models"
    "api-customer-merchant/internal/db/repositories"  // If needed for custom queries
    "github.com/rpip/paystack-go"
    "github.com/shopspring/decimal"
    "gorm.io/gorm"
    "go.uber.org/zap"
)

type PaymentService struct {
    client *paystack.Client
    conf   *config.Config
    logger *zap.Logger
    // Add repos if needed: orderRepo, etc.
}

func NewPaymentService(conf *config.Config, logger *zap.Logger) *PaymentService {
    return &PaymentService{
        client: paystack.NewClient(conf.PaystackSecretKey),
        conf:   conf,
        logger: logger,
    }
}

// InitiateTransaction: Returns auth URL for redirect
func (s *PaymentService) InitiateTransaction(ctx context.Context, order *models.Order, email string) (string, string, error) {
    req := &paystack.TransactionRequest{
        Email:      email,
        Amount:     int(order.TotalAmount.Mul(decimal.NewFromInt(100)).IntPart()),  // In kobo
        Reference:  fmt.Sprintf("order_%d", order.ID),
        Currency:   "NGN",  // Adjust as needed
        CallbackURL: s.conf.BaseURL + "/callback/paystack",  // Add callback handler if needed
    }
    resp, err := s.client.Transaction.Initialize(req)
    if err != nil {
        return "", "", err
    }
    // Update Payment with ref
    payment := &models.Payment{OrderID: order.ID, Amount: order.TotalAmount, Status: "pending", PaystackRef: resp.Data.Reference}
    if err := db.DB.Create(payment).Error; err != nil {
        return "", "", err
    }
    return resp.Data.AuthorizationURL, resp.Data.Reference, nil
}

// HandleWebhook: Verify and process
func (s *PaymentService) HandleWebhook(r *http.Request) error {
    // Verify signature
    sig := r.Header.Get("X-Paystack-Signature")
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        return err
    }
    hmac512 := hmac.New(sha512.New, []byte(s.conf.PaystackSecretKey))
    hmac512.Write(body)
    expectedSig := hex.EncodeToString(hmac512.Sum(nil))
    if sig != expectedSig {
        return errors.New("invalid signature")
    }

    var event map[string]interface{}
    if err := json.Unmarshal(body, &event); err != nil {
        return err
    }
    if event["event"] != "charge.success" {
        return nil  // Ignore other events
    }
    data := event["data"].(map[string]interface{})
    ref := data["reference"].(string)

    // Find payment/order
    var payment models.Payment
    if err := db.DB.Where("paystack_ref = ?", ref).First(&payment).Error; err != nil {
        return err
    }
    if payment.Status == "success" {
        return nil  // Idempotent
    }
    payment.Status = "success"
    if err := db.DB.Save(&payment).Error; err != nil {
        return err
    }

    var order models.Order
    if err := db.DB.Preload("OrderItems").First(&order, payment.OrderID).Error; err != nil {
        return err
    }
    order.Status = "paid"
    if err := db.DB.Save(&order).Error; err != nil {
        return err
    }

    // Calculate splits
    shares := make(map[string]decimal.Decimal)
    for _, item := range order.OrderItems {
        itemAmount := decimal.NewFromFloat(item.Price).Mul(decimal.NewFromInt(int64(item.Quantity)))  // Assume Price stored per item
        shares[item.MerchantID] = shares[item.MerchantID].Add(itemAmount)
    }
    commissionRate := decimal.NewFromFloat(s.conf.PlatformCommission)
    now := time.Now()
    for merchantID, gross := range shares {
        net := gross.Mul(decimal.NewFromFloat(1).Sub(commissionRate))
        fee := gross.Sub(net)
        holdUntil := now.Add(72 * time.Hour)  // 3 days

        split := models.OrderMerchantSplit{
            OrderID:    order.ID,
            MerchantID: merchantID,
            AmountDue:  net,
            Fee:        fee,
            Status:     "pending",
            HoldUntil:  holdUntil,
        }
        if err := db.DB.Create(&split).Error; err != nil {
            return err
        }
    }

    // Optional: Notify admin/merchant via services/notifications
    return nil
}


*/





/*

import (
	"context"
	"errors"
	"time"

	p"github.com/gray-adeyi/paystack"

	"api-customer-merchant/internal/config"
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"
	"go.uber.org/zap"
)

type PaymentService struct {
	client *p.PaystackClient
	logger *zap.Logger
	paymentRepo *repositories.PaymentRepository
	orderRepo *repositories.OrderRepository
}

func NewPaymentService(conf *config.Config, logger *zap.Logger, paymentRepo *repositories.PaymentRepository, orderRepo *repositories.OrderRepository) *PaymentService {
	client := p.NewClient(p.WithSecretKey(conf.PaystackSecretKey))
	return &PaymentService{
		client: client,
		logger: logger,
		paymentRepo: paymentRepo,
		orderRepo: orderRepo,
	}
}

func (s *PaymentService) InitializePayment(ctx context.Context, orderID uint, email string, amount float64, callbackURL string) (string, error) {
	logger := s.logger.With(zap.String("operation", "InitializePayment"), zap.Uint("order_id", orderID))

	amountInt := int(amount * 100) // Convert to kobo

	params := paystack.TransactionInitializeParams{
		Amount: amountInt,
		Email: email,
		CallbackURL: callbackURL,
		Metadata: map[string]interface{}{
			"order_id": orderID,
		},
	}

	var response paystack.Response[paystack.Transaction]
	if err := s.client.Transaction.Initialize(ctx, &params, &response); err != nil {
		logger.Error("Failed to initialize payment", zap.Error(err))
		return "", err
	}

	if !response.Status {
		return "", errors.New("payment initialization failed")
	}

	logger.Info("Payment initialized successfully", zap.String("authorization_url", response.Data.AuthorizationURL))
	return response.Data.AuthorizationURL, nil
}

func (s *PaymentService) VerifyPayment(ctx context.Context, reference string) (*models.Payment, error) {
	logger := s.logger.With(zap.String("operation", "VerifyPayment"))

	var response paystack.Response[paystack.Transaction]
	if err := s.client.Transaction.Verify(ctx, reference, &response); err != nil {
		logger.Error("Failed to verify payment", zap.Error(err))
		return nil, err
	}

	if !response.Status || response.Data.Status != "success" {
		return nil, errors.New("payment verification failed")
	}

	orderIDFloat, ok := response.Data.Metadata["order_id"].(float64)
	if !ok {
		return nil, errors.New("invalid metadata")
	}
	orderID := uint(orderIDFloat)

	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		logger.Error("Failed to find order", zap.Error(err))
		return nil, err
	}

	amount := float64(response.Data.Amount) / 100

	payment := &models.Payment{
		OrderID: order.ID,
		TransactionID: response.Data.Reference,
		Amount: amount,
		Status: models.PaymentStatusCompleted,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.paymentRepo.Create(payment); err != nil {
		logger.Error("Failed to create payment record", zap.Error(err))
		return nil, err
	}

	// Update order status
	order.Status = models.OrderStatusPaid
	if err := s.orderRepo.Update(order); err != nil {
		logger.Error("Failed to update order status", zap.Error(err))
		return nil, err
	}

	logger.Info("Payment verified successfully", zap.Uint("order_id", order.ID))
	return payment, nil
}

*/