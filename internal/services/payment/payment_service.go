package payment

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/config"
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"

	"github.com/gray-adeyi/paystack"
	m "github.com/gray-adeyi/paystack/models"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrPaymentFailed      = errors.New("payment initialization failed")
	ErrVerificationFailed = errors.New("payment verification failed")
	ErrRefundFailed       = errors.New("refund failed")
)

type PaymentService struct {
	paymentRepo *repositories.PaymentRepository
	orderRepo   *repositories.OrderRepository
	payoutRepo  *repositories.PayoutRepository
	//splitRepo    repositories.O
	merchantRepo *repositories.MerchantRepository
	//client      *paystack.Client
	config *config.Config
	logger *zap.Logger
	db     *gorm.DB
}

func NewPaymentService(
	paymentRepo *repositories.PaymentRepository,
	orderRepo *repositories.OrderRepository,
	payoutRepo *repositories.PayoutRepository,
	//splitRepo repositories.SplitRepository,
	merchantRepo *repositories.MerchantRepository,
	conf *config.Config,
	logger *zap.Logger,
) *PaymentService {

	return &PaymentService{
		paymentRepo: paymentRepo,
		orderRepo:   orderRepo,
		payoutRepo:  payoutRepo,
		//splitRepo:    splitRepo,
		merchantRepo: merchantRepo,
		//client:      client,
		config: conf,
		logger: logger,
		db:     db.DB,
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
		Status:        models.PaymentStatusPending,
		TransactionID: psResp.Data.Reference,
	}
	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		logger.Error("Failed to save payment", zap.Error(err))
		return nil, fmt.Errorf("failed to save payment: %w", err)
	}

	// Map to DTO response
	response := &dto.PaymentResponse{
		ID:               payment.ID,
		OrderID:          payment.OrderID,
		Amount:           payment.Amount.InexactFloat64(),
		Currency:         payment.Currency,
		Status:           dto.PaymentStatus(payment.Status),
		TransactionID:    payment.TransactionID,
		AuthorizationURL: psResp.Data.AuthorizationUrl, // used by frontend for redirect
		CreatedAt:        payment.CreatedAt,
		UpdatedAt:        payment.UpdatedAt,
	}

	return response, nil
}


func (s *PaymentService) VerifyPayment(ctx context.Context, reference string) (*dto.PaymentResponse, error) {
	logger := s.logger.With(zap.String("operation", "VerifyPayment"), zap.String("reference", reference))

	payment, perr := s.paymentRepo.FindByTransactionID(ctx, reference)
	if perr != nil {
		return nil, fmt.Errorf("payment not found: %w", perr)
	}

	// Prevent double processing
	if payment.Status == models.PaymentStatusCompleted {
		logger.Warn("Payment already processed")
		return s.mapPaymentToDTO(payment), nil
	}

	// psClient := paystack.NewClient(paystack.WithSecretKey(s.config.PaystackSecretKey))
	// var resp m.Response[m.Transaction]
	// err := psClient.Transactions.Verify(ctx, reference, &resp)
	resp, err := s.verifyPaystack(ctx, reference)

	if err != nil || !resp.Status || resp.Data.Status != "success" {
		logger.Error("Paystack verification failed", zap.Error(err))

		// Update payment to failed
		payment.Status = models.PaymentStatusFailed
		_ = s.paymentRepo.Update(ctx, payment)

		return nil, ErrVerificationFailed
	}
	meta, merr := normalizeMetadata(resp.Data.Metadata)
	if merr != nil {
		logger.Warn("failed to parse metadata", zap.Error(merr))
	} else if meta == nil {
		logger.Debug("metadata is empty or null")
	} else {
		logger.Debug("payment metadata", zap.Any("metadata", meta))
	}

	// Payment successful - now commit inventory and update order
	err = s.db.WithContext(context.Background()).Transaction(func(tx *gorm.DB) error {
		// Reload & lock payment
		var p models.Payment
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", payment.ID).First(&p).Error; err != nil {
			return fmt.Errorf("failed to lock payment: %w", err)
		}

		// Already completed check
		if p.Status == models.PaymentStatusCompleted {
			logger.Warn("Payment already processed (locked)")
			*payment = p
			return nil
		}

		// Update payment status
		p.Status = models.PaymentStatusCompleted
		p.UpdatedAt = time.Now()
		if err := tx.Save(&p).Error; err != nil {
			return fmt.Errorf("failed to update payment: %w", err)
		}

		// Preload order with items
		var order models.Order
		if err := tx.Preload("OrderItems").
			Where("id = ?", p.OrderID).
			First(&order).Error; err != nil {
			return fmt.Errorf("order not found: %w", err)
		}

		if len(order.OrderItems) == 0 {
			logger.Warn("order has no items", zap.Uint("order_id", order.ID))
		}

		// Collect inventory IDs and quantities for batch update
		type invUpdate struct {
			ID       string
			Quantity int
		}
		var updates []invUpdate
		invIDs := make([]string, 0, len(order.OrderItems))

		for _, item := range order.OrderItems {
			var inv models.Inventory
			q := tx.Where("merchant_id = ?", item.MerchantID)
			if item.VariantID != nil && *item.VariantID != "" {
				q = q.Where("variant_id = ?", *item.VariantID)
			} else {
				q = q.Where("product_id = ?", item.ProductID)
			}

			// Lock inventory row
			if err := q.Clauses(clause.Locking{Strength: "UPDATE"}).First(&inv).Error; err != nil {
				return fmt.Errorf("inventory not found for item %v: %w", item.ID, err)
			}

			updates = append(updates, invUpdate{ID: inv.ID, Quantity: item.Quantity})
			invIDs = append(invIDs, inv.ID)
		}

		// Batch update inventories
		if len(updates) > 0 {
			var casesQty, casesRes []string
			for _, u := range updates {
				casesQty = append(casesQty, fmt.Sprintf("WHEN id = '%s' THEN GREATEST(quantity - %d, 0)", u.ID, u.Quantity))
				casesRes = append(casesRes, fmt.Sprintf("WHEN id = '%s' THEN GREATEST(reserved_quantity - %d, 0)", u.ID, u.Quantity))
			}

			sql := fmt.Sprintf(`
				UPDATE inventories
				SET quantity = CASE %s END,
					reserved_quantity = CASE %s END,
					updated_at = NOW()
				WHERE id IN (?);
			`, strings.Join(casesQty, " "), strings.Join(casesRes, " "))

			if err := tx.Exec(sql, invIDs).Error; err != nil {
				return fmt.Errorf("failed to batch update inventories: %w", err)
			}
		}

		// Update order status to Paid
		order.Status = models.OrderStatusProcessing
		if err := tx.Save(&order).Error; err != nil {
			return fmt.Errorf("failed to update order status: %w", err)
		}

		// Update merchant splits to processing
		if err := tx.Model(&models.OrderMerchantSplit{}).
			Where("order_id = ? AND status = ?", order.ID, "pending").
			Update("status", "processing").Error; err != nil {
			return fmt.Errorf("failed to update merchant splits: %w", err)
		}

		// Clear cart items using join (fix user_id issue)
		if err := tx.Exec(`
			DELETE FROM cart_items
			USING carts
			WHERE cart_items.cart_id = carts.id
				AND carts.user_id = ?
				AND carts.status = ?
		`, order.UserID, models.CartStatusActive).Error; err != nil {
			return fmt.Errorf("failed to clear cart items: %w", err)
		}

		// Update cart status to Converted
		if err := tx.Model(&models.Cart{}).
			Where("user_id = ? AND status = ?", order.UserID, models.CartStatusActive).
			UpdateColumn("status", models.CartStatusConverted).Error; err != nil {
			return fmt.Errorf("failed to update cart status: %w", err)
		}

		// Reflect updated payment
		*payment = p
		return nil
	})

	if err != nil {
		logger.Error("Failed to commit payment transaction", zap.Error(err))
		return nil, fmt.Errorf("failed to commit payment: %w", err)
	}

	logger.Info("Payment verified and committed",
		zap.Uint("payment_id", payment.ID),
		zap.Uint("order_id", payment.OrderID),
	)

	return s.mapPaymentToDTO(payment), nil
}

func (s *PaymentService) mapPaymentToDTO(payment *models.Payment) *dto.PaymentResponse {
	return &dto.PaymentResponse{
		ID:            payment.ID,
		OrderID:       payment.OrderID,
		Amount:        payment.Amount.InexactFloat64(),
		Currency:      payment.Currency,
		Status:        dto.PaymentStatus(payment.Status),
		TransactionID: payment.TransactionID,
		CreatedAt:     payment.CreatedAt,
		UpdatedAt:     payment.UpdatedAt,
	}
}

func (s *PaymentService) handleChargeSuccess(ctx context.Context, reference string) (*dto.PaymentResponse, error) {
	logger := s.logger.With(zap.String("operation", "VerifyPayment"), zap.String("reference", reference))

	payment, perr := s.paymentRepo.FindByTransactionID(ctx, reference)
	if perr != nil {
		return nil, fmt.Errorf("payment not found: %w", perr)
	}

	// Prevent double processing
	if payment.Status == models.PaymentStatusCompleted {
		logger.Warn("Payment already processed")
		return s.mapPaymentToDTO(payment), nil
	}

	// Payment successful - now commit inventory and update order
	err := s.db.WithContext(context.Background()).Transaction(func(tx *gorm.DB) error {
		// Reload & lock payment
		var p models.Payment
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", payment.ID).First(&p).Error; err != nil {
			return fmt.Errorf("failed to lock payment: %w", err)
		}

		// Already completed check
		if p.Status == models.PaymentStatusCompleted {
			logger.Warn("Payment already processed (locked)")
			*payment = p
			return nil
		}

		// Update payment status
		p.Status = models.PaymentStatusCompleted
		p.UpdatedAt = time.Now()
		if err := tx.Save(&p).Error; err != nil {
			return fmt.Errorf("failed to update payment: %w", err)
		}

		// Preload order with items
		var order models.Order
		if err := tx.Preload("OrderItems").
			Where("id = ?", p.OrderID).
			First(&order).Error; err != nil {
			return fmt.Errorf("order not found: %w", err)
		}

		if len(order.OrderItems) == 0 {
			logger.Warn("order has no items", zap.Uint("order_id", order.ID))
		}

		// Collect inventory IDs and quantities for batch update
		type invUpdate struct {
			ID       string
			Quantity int
		}
		var updates []invUpdate
		invIDs := make([]string, 0, len(order.OrderItems))

		for _, item := range order.OrderItems {
			var inv models.Inventory
			q := tx.Where("merchant_id = ?", item.MerchantID)
			if item.VariantID != nil && *item.VariantID != "" {
				q = q.Where("variant_id = ?", *item.VariantID)
			} else {
				q = q.Where("product_id = ?", item.ProductID)
			}

			// Lock inventory row
			if err := q.Clauses(clause.Locking{Strength: "UPDATE"}).First(&inv).Error; err != nil {
				return fmt.Errorf("inventory not found for item %v: %w", item.ID, err)
			}

			updates = append(updates, invUpdate{ID: inv.ID, Quantity: item.Quantity})
			invIDs = append(invIDs, inv.ID)
		}

		// Batch update inventories
		if len(updates) > 0 {
			var casesQty, casesRes []string
			for _, u := range updates {
				casesQty = append(casesQty, fmt.Sprintf("WHEN id = '%s' THEN GREATEST(quantity - %d, 0)", u.ID, u.Quantity))
				casesRes = append(casesRes, fmt.Sprintf("WHEN id = '%s' THEN GREATEST(reserved_quantity - %d, 0)", u.ID, u.Quantity))
			}

			sql := fmt.Sprintf(`
				UPDATE inventories
				SET quantity = CASE %s END,
					reserved_quantity = CASE %s END,
					updated_at = NOW()
				WHERE id IN (?);
			`, strings.Join(casesQty, " "), strings.Join(casesRes, " "))

			if err := tx.Exec(sql, invIDs).Error; err != nil {
				return fmt.Errorf("failed to batch update inventories: %w", err)
			}
		}

		// Update order status to Paid
		order.Status = models.OrderStatusProcessing
		if err := tx.Save(&order).Error; err != nil {
			return fmt.Errorf("failed to update order status: %w", err)
		}

		// Update merchant splits to processing
		if err := tx.Model(&models.OrderMerchantSplit{}).
			Where("order_id = ? AND status = ?", order.ID, "pending").
			Update("status", "processing").Error; err != nil {
			return fmt.Errorf("failed to update merchant splits: %w", err)
		}

		// Clear cart items using join (fix user_id issue)
		if err := tx.Exec(`
			DELETE FROM cart_items
			USING carts
			WHERE cart_items.cart_id = carts.id
				AND carts.user_id = ?
				AND carts.status = ?
		`, order.UserID, models.CartStatusActive).Error; err != nil {
			return fmt.Errorf("failed to clear cart items: %w", err)
		}

		// Update cart status to Converted
		if err := tx.Model(&models.Cart{}).
			Where("user_id = ? AND status = ?", order.UserID, models.CartStatusActive).
			UpdateColumn("status", models.CartStatusConverted).Error; err != nil {
			return fmt.Errorf("failed to update cart status: %w", err)
		}

		// Reflect updated payment
		*payment = p
		return nil
	})

	if err != nil {
		logger.Error("Failed to commit payment transaction", zap.Error(err))
		return nil, fmt.Errorf("failed to commit payment: %w", err)
	}

	logger.Info("Payment verified and committed",
		zap.Uint("payment_id", payment.ID),
		zap.Uint("order_id", payment.OrderID),
	)

	return s.mapPaymentToDTO(payment), nil
}

// GetPaymentByOrderID retrieves a payment by order ID
func (s *PaymentService) GetPaymentByOrderID(ctx context.Context, orderID uint) (*models.Payment, error) {
	if orderID == 0 {
		return nil, errors.New("invalid order ID")
	}
	return s.paymentRepo.FindByOrderID(ctx, orderID)
}

// GetPaymentsByUserID retrieves all payments for a user
func (s *PaymentService) GetPaymentsByUserID(ctx context.Context, userID uint) ([]models.Payment, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}
	return s.paymentRepo.FindByUserID(ctx, userID)
}

// UpdatePaymentStatus updates the status of a payment
func (s *PaymentService) UpdatePaymentStatus(ctx context.Context, paymentID uint, status string) (*models.Payment, error) {
	if paymentID == 0 {
		return nil, errors.New("invalid payment ID")
	}
	if err := models.PaymentStatus(status).Valid(); err != nil {
		return nil, err
	}

	payment, err := s.paymentRepo.FindByID(ctx, paymentID)
	if err != nil {
		return nil, err
	}

	payment.Status = models.PaymentStatus(status)
	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		return nil, err
	}

	return s.paymentRepo.FindByID(ctx, paymentID)
}

// HandleWebhook processes the forwarded Paystack event
func (s *PaymentService) HandleWebhook(ctx context.Context, event map[string]interface{}) error {
	eventType, ok := event["event"].(string)
	if !ok {
		return errors.New("invalid event type")
	}

	data, ok := event["data"].(map[string]interface{})
	if !ok {
		return errors.New("invalid event data")
	}
	status, _ := data["status"].(string)
	validStatuses := map[string]bool{"success": true, "failed": true, "abandoned": true}
	if eventType == "charge.success" && !validStatuses[status] {
		s.logger.Warn("Invalid payment status in webhook", zap.String("status", status))
		return fmt.Errorf("invalid payment status: %s", status)
	}

	s.logger.Info("Received Paystack webhook", zap.String("event", eventType))

	switch eventType {
	case "transfer.success":
		return s.handleTransferSuccess(ctx, data)
	case "transfer.failed", "transfer.reversed":
		return s.handleTransferFailure(ctx, data)
	case "charge.success":
		// Handle if needed for order processing
		s.logger.Info("Charge success event received", zap.Any("reference", data["reference"]))

		reference, ok := data["reference"].(string)
		if !ok {
			return errors.New("invalid reference type")
		}

		_, err := s.handleChargeSuccess(ctx, reference)
		return err
	default:
		s.logger.Info("Unhandled webhook event", zap.String("event", eventType))
		return nil
	}
}

// handleTransferSuccess handles transfer.success event
func (s *PaymentService) handleTransferSuccess(ctx context.Context, data map[string]interface{}) error {
	transferCode, ok := data["recipient_code"].(string)
	if !ok {
		return errors.New("invalid transfer_code")
	}

	payout, err := s.payoutRepo.FindByPaystackTransferID(ctx, transferCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("No payout found for transfer code", zap.String("transfer_code", transferCode))
			return nil
		}
		return err
	}

	payout.Status = "completed"
	if err := s.payoutRepo.Update(ctx, payout); err != nil {
		return err
	}

	// Mark related splits as paid
	// if err := s.splitRepo.UpdateStatusByMerchantAndStatus(ctx, payout.MerchantID, "processing", "paid"); err != nil {
	// 	return err
	// }

	// Update merchant totals
	merchant, err := s.merchantRepo.GetByMerchantID(ctx, payout.MerchantID)
	if err != nil {
		return err
	}

	merchant.TotalPayouts = merchant.TotalPayouts + payout.Amount
	merchant.LastPayoutDate = func() *time.Time { t := time.Now(); return &t }() // if err := s.merchantRepo.Update(ctx, merchant); err != nil {
	// 	return err
	// }

	s.logger.Info("Payout completed successfully", zap.String("payout_id", payout.ID))
	return nil
}

// handleTransferFailure handles transfer.failed or reversed
func (s *PaymentService) handleTransferFailure(ctx context.Context, data map[string]interface{}) error {
	transferCode, ok := data["transfer_code"].(string)
	if !ok {
		return errors.New("invalid transfer_code")
	}

	payout, err := s.payoutRepo.FindByPaystackTransferID(ctx, transferCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("No payout found for transfer code", zap.String("transfer_code", transferCode))
			return nil
		}
		return err
	}

	payout.Status = "failed"
	if err := s.payoutRepo.Update(ctx, payout); err != nil {
		return err
	}

	// Reset splits to payout_requested
	// if err := s.splitRepo.UpdateStatusByMerchantAndStatus(ctx, payout.MerchantID, "processing", "payout_requested"); err != nil {
	// 	return err
	// }

	// Send notification (implement if needed)
	reason, _ := data["reason"].(string)
	// merchant, err := s.merchantRepo.FindByID(ctx, payout.MerchantID)
	// if err == nil {
	// 	// Call sendPayoutFailedEmail(merchant.WorkEmail, merchant.StoreName, payout.Amount.InexactFloat64(), reason)
	// 	// Stub or implement the email function
	// }

	s.logger.Error("Payout failed", zap.String("payout_id", payout.ID), zap.String("reason", reason))
	return nil
}
