package payment


import (
	"context"
	"errors"
	"fmt"
	"time"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/config"
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"

	"github.com/gray-adeyi/paystack"
	m "github.com/gray-adeyi/paystack/models"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

var (
	ErrPaymentFailed     = errors.New("payment initialization failed")
	ErrVerificationFailed = errors.New("payment verification failed")
	ErrRefundFailed      = errors.New("refund failed")
)


type PaymentService struct {
	paymentRepo repositories.PaymentRepository
	orderRepo   repositories.OrderRepository
	//client      *paystack.Client
	config      *config.Config
	logger      *zap.Logger
}

func NewPaymentService(
	paymentRepo repositories.PaymentRepository,
	orderRepo repositories.OrderRepository,
	conf *config.Config,
	logger *zap.Logger,
) *PaymentService {
	//client := paystack.NewClient(paystack.WithSecretKey(conf.PaystackSecretKey))
	return &PaymentService{
		paymentRepo: paymentRepo,
		orderRepo:   orderRepo,
		//client:      client,
		config:      conf,
		logger:      logger,
	}
}

/*
func (s *PaymentService) InitializeCheckout(ctx context.Context, req dto.InitializePaymentRequest) (*dto.PaymentResponse, error) {
	logger := s.logger.With(zap.Uint("order_id", req.OrderID))

	// Validate and fetch order
	order, err := s.orderRepo.FindByID(ctx,req.OrderID)
	if err != nil {
		logger.Error("Order not found", zap.Error(err))
		return nil, fmt.Errorf("order not found: %w", err)
	}
	amountInKobo := order.SubTotal.Mul(decimal.NewFromInt(100)).String()
	psClient := paystack.NewClient(paystack.WithSecretKey(s.config.PaystackSecretKey))
	// Verify amount matches order total
	expectedTotal := order.SubTotal.InexactFloat64() // Assuming decimal.Decimal
	if req.Amount != expectedTotal {
		logger.Error("Amount mismatch", zap.Float64("expected", expectedTotal), zap.Float64("got", req.Amount))
		return nil, fmt.Errorf("amount mismatch: expected %v, got %v", expectedTotal, req.Amount)
	}

	// Convert amount to kobo for Paystack
	//amountKobo := int(req.Amount * 100)

	// Initialize Paystack transaction
	initReq := &paystack.InitializeTransactionRequest{
		Amount:    amountInKobo,
		Email:     req.Email,
		Currency:  req.Currency,
		Reference: fmt.Sprintf("order_%d", order.ID),
	}
	resp, err := psClient.Transactions.Initialize(context.TODO(),initReq)
	if err != nil {
		logger.Error("Paystack init failed", zap.Error(err))
		return nil, fmt.Errorf("paystack init failed: %w", err)
	}

	// Save payment
	payment := &models.Payment{
		OrderID:       order.ID,
		Amount:        decimal.NewFromFloat(req.Amount),
		Currency:      req.Currency,
		Status:        "pending",
		TransactionID: resp.Data.Reference,
	}
	if err := s.paymentRepo.Create(payment); err != nil {
		logger.Error("Failed to save payment", zap.Error(err))
		return nil, fmt.Errorf("failed to save payment: %w", err)
	}

	// Manual mapping
	response := &dto.PaymentResponse{
		ID:             payment.ID,
		OrderID:        payment.OrderID,
		Amount:         payment.Amount.InexactFloat64(),
		Currency:       payment.Currency,
		Status:         payment.Status,
		TransactionID:  payment.TransactionID,
		AuthorizationURL: resp.Data.AuthorizationURL, // For frontend redirect
		CreatedAt:      payment.CreatedAt,
		UpdatedAt:      payment.UpdatedAt,
	}
	return response, nil
}
*/









func (s *PaymentService) InitializeCheckout(ctx context.Context, req dto.InitializePaymentRequest) (*dto.PaymentResponse, error) {
	logger := s.logger.With(zap.Uint("order_id", req.OrderID))

	// Validate and fetch order
	order, err := s.orderRepo.FindByID(ctx, req.OrderID)
	if err != nil {
		logger.Error("Order not found", zap.Error(err))
		return nil, fmt.Errorf("order not found: %w", err)
	}

	// Compute expected amount in kobo (integer) from order subtotal (decimal.Decimal)
	amountKoboFromOrderInt64 := order.SubTotal.Mul(decimal.NewFromInt(100)).IntPart()
	amountKobo := int(amountKoboFromOrderInt64) // paystack client expects int (kobo)
	// Compute requested amount in kobo and compare as integers to avoid float equality issues
	reqAmountKoboInt64 := decimal.NewFromFloat(req.Amount).Mul(decimal.NewFromInt(100)).IntPart()
	reqAmountKobo := int(reqAmountKoboInt64)

	if reqAmountKobo != amountKobo {
		logger.Error("Amount mismatch", zap.Int("expected_kobo", amountKobo), zap.Int("got_kobo", reqAmountKobo))
		return nil, fmt.Errorf("amount mismatch: expected %d kobo, got %d kobo", amountKobo, reqAmountKobo)
	}

	// Create Paystack client
	psClient := paystack.NewClient(paystack.WithSecretKey(s.config.PaystackSecretKey))

	// Call Transactions.Initialize(amount int, email string, response any, optionalPayloads ...)
	var psResp m.Response[m.InitTransaction]
	// pass currency as optional payload if provided
	var initErr error
	if req.Currency != "" {
		initErr = psClient.Transactions.Initialize(ctx, amountKobo, req.Email, &psResp, paystack.WithOptionalPayload("currency", req.Currency))
	} else {
		initErr = psClient.Transactions.Initialize(ctx, amountKobo, req.Email, &psResp)
	}
	if initErr != nil {
		logger.Error("Paystack initialize transaction failed", zap.Error(initErr))
		return nil, fmt.Errorf("paystack initialize failed: %w", initErr)
	}

	// Ensure response data exists
	if psResp.Data.Reference == "" {
		logger.Error("Paystack initialize returned empty reference", zap.Any("response", psResp))
		return nil, fmt.Errorf("paystack initialize returned empty reference")
	}

	// Save payment. Use order.SubTotal (the canonical amount) to avoid any float round issues.
	payment := &models.Payment{
		OrderID:       order.ID,
		Amount:        order.SubTotal, // decimal.Decimal
		Currency:      req.Currency,
		Status:        "pending",
		TransactionID: psResp.Data.Reference,
	}
	if err := s.paymentRepo.Create(ctx,payment); err != nil {
		logger.Error("Failed to save payment", zap.Error(err))
		return nil, fmt.Errorf("failed to save payment: %w", err)
	}

	// Map to DTO response
	response := &dto.PaymentResponse{
		ID:               payment.ID,
		OrderID:          payment.OrderID,
		Amount:           payment.Amount.InexactFloat64(),
		Currency:         payment.Currency,
		Status:          string(payment.Status),
		TransactionID:    payment.TransactionID,
		AuthorizationURL: psResp.Data.AuthorizationUrl, // used by frontend for redirect
		CreatedAt:        payment.CreatedAt,
		UpdatedAt:        payment.UpdatedAt,
	}

	return response, nil
}








/*
func (s *PaymentService) VerifyPayment(ctx context.Context, reference string) (*dto.PaymentResponse, error) {
	logger := s.logger.With(zap.String("reference", reference))
	


	// Verify with Paystack
	psClient := paystack.NewClient(paystack.WithSecretKey(s.config.PaystackSecretKey))
	var psResp m.Response[m.Transaction]
	 err := psClient.Transactions.Verify(ctx, reference, &psResp)
	if err != nil {
		logger.Error("Paystack verify failed", zap.Error(err))
		return nil, fmt.Errorf("paystack verify failed: %w", err)
	}
	if psResp.Data.Status != "success" {
		logger.Warn("Payment not successful", zap.String("status", string(psResp.Data.Status)))
		return nil, fmt.Errorf("payment not successful: status %s",string(psResp.Data.Status))
	}

	// Fetch and update payment
	payment, err := s.paymentRepo.FindByTransactionID(ctx, reference)
	if err != nil {
		logger.Error("Payment not found", zap.Error(err))
		return nil, fmt.Errorf("payment not found: %w", err)
	}
	payment.Status = "success"
	if err := s.paymentRepo.Update(payment); err != nil {
		logger.Error("Failed to update payment", zap.Error(err))
		return nil, fmt.Errorf("failed to update payment: %w", err)
	}

	// Update order status
	order, err := s.orderRepo.FindByID(ctx ,payment.OrderID)
	if err != nil {
		logger.Error("Order not found", zap.Error(err))
		return nil, fmt.Errorf("order not found: %w", err)
	}
	order.Status = "paid"
	if err := s.orderRepo.Update(order); err != nil {
		logger.Error("Failed to update order", zap.Error(err))
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	// Manual mapping
	response := &dto.PaymentResponse{
		ID:            payment.ID,
		OrderID:       payment.OrderID,
		Amount:        payment.Amount.InexactFloat64(),
		Currency:      payment.Currency,
		Status:       string(payment.Status),
		TransactionID: payment.TransactionID,
		CreatedAt:     payment.CreatedAt,
		UpdatedAt:     payment.UpdatedAt,
	}
	return response, nil
}
*/


func (s *PaymentService) VerifyPayment(ctx context.Context, reference string) (*dto.PaymentResponse, error) {
	logger := s.logger.With(zap.String("operation", "VerifyPayment"), zap.String("reference", reference))

	payment, perr := s.paymentRepo.FindByTransactionID(ctx, reference)
	if perr != nil {
		return nil, fmt.Errorf("payment not found: %w", perr)
	}

	psClient := paystack.NewClient(paystack.WithSecretKey(s.config.PaystackSecretKey))
	var resp m.Response[m.Transaction]
	 err := psClient.Transactions.Verify(ctx, reference, &resp)
	if err != nil || !resp.Status  || resp.Data.Status != "success" {
		logger.Error("Paystack verification failed", zap.Error(err))
		// Update to failed
		payment.Status = models.PaymentStatusFailed
		s.paymentRepo.Update(ctx, payment)
		return nil, ErrVerificationFailed
	}

	// Update success
	payment.Status = models.PaymentStatusCompleted
	payment.UpdatedAt = time.Now()
	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		return nil, err
	}

	// Update order status
	order, err := s.orderRepo.FindByID(ctx, payment.OrderID)
	if err == nil {
		order.Status = models.OrderStatusCompleted
		s.orderRepo.Update(ctx, order)
	}

	logger.Info("Payment verified", zap.Uint("payment_id", payment.ID))
	response := &dto.PaymentResponse{
		ID:            payment.ID,
		OrderID:       payment.OrderID,
		Amount:        payment.Amount.InexactFloat64(),
		Currency:      payment.Currency,
		Status:       string(payment.Status),
		TransactionID: payment.TransactionID,
		CreatedAt:     payment.CreatedAt,
		UpdatedAt:     payment.UpdatedAt,
	}
	return response, nil
}



// GetPaymentByOrderID retrieves a payment by order ID
func (s *PaymentService) GetPaymentByOrderID(ctx context.Context, orderID uint) (*models.Payment, error) {
	if orderID == 0 {
		return nil, errors.New("invalid order ID")
	}
	return s.paymentRepo.FindByOrderID(ctx,orderID)
}

// GetPaymentsByUserID retrieves all payments for a user
func (s *PaymentService) GetPaymentsByUserID(ctx context.Context,userID uint) ([]models.Payment, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}
	return s.paymentRepo.FindByUserID(ctx,userID)
}

// UpdatePaymentStatus updates the status of a payment
func (s *PaymentService) UpdatePaymentStatus(ctx context.Context,paymentID uint, status string) (*models.Payment, error) {
	if paymentID == 0 {
		return nil, errors.New("invalid payment ID")
	}
	if err := models.PaymentStatus(status).Valid(); err != nil {
		return nil, err
	}

	payment, err := s.paymentRepo.FindByID(ctx ,paymentID)
	if err != nil {
		return nil, err
	}

	payment.Status = models.PaymentStatus(status)
	if err := s.paymentRepo.Update(ctx ,payment); err != nil {
		return nil, err
	}

	return s.paymentRepo.FindByID(ctx ,paymentID)
}




func (s *PaymentService) HandleWebhook() {
	// logger := s.logger.With(zap.String("operation", "HandleWebhook"))

	// var event paystack.Event
	// if err := c.ShouldBindJSON(&event); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
	// 	return
	// }

	// // Verify signature
	// sig := c.GetHeader("x-paystack-signature")
	// if !paystack.VerifySignature([]byte(event.Raw), sig, s.conf.PaystackSecretKey) {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid signature"})
	// 	return
	// }

	// // Handle charge.success
	// if event.Event == "charge.success" {
	// 	var data paystack.TransactionData
	// 	if err := json.Unmarshal(event.Data, &data); err != nil {
	// 		logger.Error("Unmarshal webhook data failed", zap.Error(err))
	// 		c.Status(http.StatusBadRequest)
	// 		return
	// 	}
	// 	paymentResp, err := s.VerifyPayment(c.Request.Context(), data.Reference)
	// 	if err != nil {
	// 		logger.Error("Webhook verification failed", zap.Error(err))
	// 		c.Status(http.StatusInternalServerError)
	// 		return
	// 	}
	// 	// Trigger notifications, payouts, etc.
	// }

	// c.JSON(http.StatusOK, gin.H{"status": "success"})
}