package order

import (
	"context"
	"errors"
	"fmt"
	//"strings"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/api/helpers"
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"
	"api-customer-merchant/internal/services/email"
	"api-customer-merchant/internal/services/payment"

	//"go.uber.org/zap"
	//"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// OrderService provides business logic for handling orders.
type OrderService struct {
	orderRepo      *repositories.OrderRepository
	orderItemRepo  *repositories.OrderItemRepository
	cartRepo       *repositories.CartRepository
	cartItemRepo   *repositories.CartItemRepository
	productRepo    *repositories.ProductRepository
	inventoryRepo  *repositories.InventoryRepository
	userRepo       *repositories.UserRepository // ADD THIS
	paymentService *payment.PaymentService
	emailService   *email.EmailService
	logger         *zap.Logger
	//validator   *validator.Validate
	db *gorm.DB
}

// NewOrderService creates a new instance of OrderService.
func NewOrderService(
	orderRepo *repositories.OrderRepository,
	orderItemRepo *repositories.OrderItemRepository,
	cartRepo *repositories.CartRepository,
	cartItemRepo *repositories.CartItemRepository,
	productRepo *repositories.ProductRepository,
	inventoryRepo *repositories.InventoryRepository,
	userRepo *repositories.UserRepository, // ADD THIS
	paymentService *payment.PaymentService,
	emailService *email.EmailService,
	logger *zap.Logger,
) *OrderService {
	return &OrderService{
		orderRepo:      orderRepo,
		orderItemRepo:  orderItemRepo,
		cartRepo:       cartRepo,
		cartItemRepo:   cartItemRepo,
		productRepo:    productRepo,
		inventoryRepo:  inventoryRepo,
		userRepo:       userRepo,
		paymentService: paymentService,
		emailService:   emailService,
		logger:         logger,
		db:             db.DB,
	}
}

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrInvalidOrderStatus = errors.New("order cannot be cancelled")
	ErrUnauthorizedOrder  = errors.New("unauthorized to cancel this order")
	ErrRefundFailed       = errors.New("failed to initiate refund")
	ErrNotificationFailed = errors.New("failed to send notification")
)

// CreateOrder converts a user's active cart into an order.
// It performs several actions within a single database transaction:
// 1. Finds the user's active cart.
// 2. Validates that the cart is not empty.
// 3. For each item in the cart, it moves the reserved stock to committed stock.
// 4. Creates an Order record.
// 5. Creates OrderItem records corresponding to the CartItems.
// 6. Deletes the cart items.
// 7. Updates the cart status to 'Converted'.
// 8. Returns a DTO representing the newly created order.

// func (s *OrderService) CreateOrder(ctx context.Context, userID uint) (*dto.OrderResponse, error) {
//     if userID == 0 {
//         return nil, errors.New("invalid user ID")
//     }

//     cart, err := s.cartRepo.FindActiveCart(ctx, userID)
//     if err != nil {
//         if errors.Is(err, repositories.ErrCartNotFound) {
//             return nil, errors.New("no active cart found")
//         }
//         return nil, err
//     }

//     if len(cart.CartItems) == 0 {
//         return nil, errors.New("cart is empty")
//     }

//     var newOrder *models.Order
//     var totalAmount decimal.Decimal
// 	type inventoryUpdate struct {
// 		ID       string
// 		Quantity int
// 	}
// 	var updates []inventoryUpdate

//     // Use a transaction to ensure atomicity
//     err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

// 		defer func() {
// 			if r := recover(); r != nil {
// 				s.logger.Error("Transaction panic", zap.Any("panic", r))
// 			}
// 		}()
//         // Calculate total and create order items
//         var orderItems []models.OrderItem
//         for _, item := range cart.CartItems {
//             // Determine price: variant if present, else product
//             price := item.Product.FinalPrice
//             if item.VariantID != nil && item.Variant != nil {
//                 price = item.Variant.FinalPrice
//             }
//             priceFloat := price.InexactFloat64()
//             subtotal := decimal.NewFromInt(int64(item.Quantity)).Mul(price)
//             totalAmount = totalAmount.Add(subtotal)

//             orderItem := models.OrderItem{
//                 ProductID:         item.ProductID,
//                 VariantID:         item.VariantID, // Set variant if present
//                 MerchantID:        item.MerchantID,
//                 Quantity:          item.Quantity,
//                 Price:             priceFloat,
//                 FulfillmentStatus: models.FulfillmentStatusNew,
//             }

//             orderItems = append(orderItems, orderItem)

//             // Commit inventory: subtract from quantity and reserved
//             // Use inventory repo for consistency, but since tx, use tx directly
//             inventoryQuery := "merchant_id = ?"
// 			args := []interface{}{item.MerchantID}
// 			if item.VariantID != nil {
// 				inventoryQuery += " AND variant_id = ?"
// 				args = append(args, *item.VariantID)
// 			} else {
// 				inventoryQuery += " AND product_id = ?"
// 				args = append(args, item.ProductID)
// 			}

// 			var inv models.Inventory
// 			if err := tx.Where(inventoryQuery, args...).First(&inv).Error; err != nil {
// 				return fmt.Errorf("inventory not found: %w", err)
// 			}
// 			updates = append(updates, inventoryUpdate{ID: inv.ID, Quantity: item.Quantity})
// 			}

// 			// Batch update using CASE WHEN
// 			if len(updates) > 0 {
// 			ids := make([]string, len(updates))
// 			casesQty := make([]string, len(updates))
// 			casesRes := make([]string, len(updates))

// 			for i, u := range updates {
// 				ids[i] = u.ID
// 				casesQty[i] = fmt.Sprintf("WHEN id = '%s' THEN quantity - %d", u.ID, u.Quantity)
// 				casesRes[i] = fmt.Sprintf("WHEN id = '%s' THEN reserved_quantity - %d", u.ID, u.Quantity)
// 			}

// 			sql := fmt.Sprintf(`
// 				UPDATE inventories
// 				SET quantity = CASE %s END,
// 					reserved_quantity = CASE %s END
// 				WHERE id IN (?)
// 			`, strings.Join(casesQty, " "), strings.Join(casesRes, " "))

// 			if err := tx.Exec(sql, ids).Error; err != nil {
// 				return fmt.Errorf("failed to batch commit stock: %w", err)
// 			}
// 				}

//         // Create the order
//         newOrder = &models.Order{
//             UserID:      userID,
//             SubTotal:    totalAmount, // Assuming SubTotal is before tax/shipping
//             TotalAmount: totalAmount, // Update if tax/shipping added later
//             Status:      models.OrderStatusPending,
//             Currency:    "NGN", // Default; make configurable
//         }
//         if err := tx.Create(newOrder).Error; err != nil {
//             return fmt.Errorf("failed to create order: %w", err)
//         }

//         // Associate and create order items
//         for i := range orderItems {
//             orderItems[i].OrderID = newOrder.ID
//         }
//         if err := tx.Create(&orderItems).Error; err != nil {
//             return fmt.Errorf("failed to create order items: %w", err)
//         }

//         // Reload order with items for response
//         if err := tx.Preload("OrderItems").First(newOrder, "id = ?", newOrder.ID).Error; err != nil {
//             return fmt.Errorf("failed to reload order: %w", err)
//         }

//         // Clear cart items
//         if err := tx.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
//             return fmt.Errorf("failed to clear cart items: %w", err)
//         }

//         // Mark cart as converted
//         cart.Status = models.CartStatusConverted
//         if err := tx.Save(cart).Error; err != nil {
//             return fmt.Errorf("failed to update cart status: %w", err)
//         }

//         return nil
//     })

//     if err != nil {
//         //s.logger.Error("Transaction failed", zap.Error(err))
//         return nil, fmt.Errorf("transaction failed: %w", err)
//     }

//     // Convert to DTO for response
// 	response := helpers.ToOrderResponse(newOrder)

// 	//  user, _ := userRepo.FindByID(userID)  // Assume userRepo injected
//     // authURL, ref, err := paymentService.InitiateTransaction(ctx, order, user.Email)
//     // if err != nil {
//     //     return nil, err
//     // }
//     // Return order with authURL in response

//     return response, nil
// }

func (s *OrderService) CreateOrder(ctx context.Context, userID uint) (*dto.OrderResponse, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}

	// Fetch user for payment initialization
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to fetch user", zap.Uint("user_id", userID), zap.Error(err))
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	cart, err := s.cartRepo.FindActiveCart(ctx, userID)
	if err != nil {
		if errors.Is(err, repositories.ErrCartNotFound) {
			return nil, errors.New("no active cart found")
		}
		return nil, err
	}

	if len(cart.CartItems) == 0 {
		return nil, errors.New("cart is empty")
	}

	var newOrder *models.Order
	var totalAmount decimal.Decimal

	// Use a transaction to ensure atomicity
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		defer func() {
			if r := recover(); r != nil {
				s.logger.Error("Transaction panic", zap.Any("panic", r))
			}
		}()

		// Calculate total and create order items (DO NOT commit inventory yet)
		var orderItems []models.OrderItem

		for _, item := range cart.CartItems {
			// Determine price: variant if present, else product
			price := item.Product.FinalPrice
			if item.VariantID != nil && item.Variant != nil {
				price = item.Variant.FinalPrice
			}
			priceFloat := price.InexactFloat64()
			subtotal := decimal.NewFromInt(int64(item.Quantity)).Mul(price)
			totalAmount = totalAmount.Add(subtotal)

			orderItem := models.OrderItem{
				ProductID:         item.ProductID,
				VariantID:         item.VariantID,
				MerchantID:        item.MerchantID,
				Quantity:          item.Quantity,
				Price:             priceFloat,
				FulfillmentStatus: models.FulfillmentStatusProcessing,
			}

			orderItems = append(orderItems, orderItem)

			// IMPORTANT: DO NOT commit inventory here
			// Inventory is already reserved in cart
			// It will only be committed after successful payment
		}

		// Create the order
		newOrder = &models.Order{
			UserID:      userID,
			SubTotal:    totalAmount,
			TotalAmount: totalAmount, // Add tax/shipping calculation if needed
			Status:      models.OrderStatusPending,
			Currency:    "NGN",
		}
		if err := tx.Create(newOrder).Error; err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		// Associate and create order items
		for i := range orderItems {
			orderItems[i].OrderID = newOrder.ID
		}
		if err := tx.Create(&orderItems).Error; err != nil {
			return fmt.Errorf("failed to create order items: %w", err)
		}

		// Reload order with items and user for response
		if err := tx.Preload("OrderItems.Product.Media").
			Preload("OrderItems.Product").
			Preload("User.Addresses").
			First(newOrder, "id = ?", newOrder.ID).Error; err != nil {
			return fmt.Errorf("failed to reload order: %w", err)
		}

		// DO NOT clear cart items yet - keep them until payment succeeds
		// DO NOT mark cart as converted yet

		return nil
	})

	if err != nil {
		s.logger.Error("Transaction failed", zap.Uint("user_id", userID), zap.Error(err))
		return nil, fmt.Errorf("transaction failed: %w", err)
	}

	// Initialize payment with Paystack
	paymentReq := dto.InitializePaymentRequest{
		OrderID:  newOrder.ID,
		Amount:   totalAmount.InexactFloat64(),
		Email:    user.Email,
		Currency: "NGN",
	}

	paymentResp, err := s.paymentService.InitializeCheckout(ctx, paymentReq)
	if err != nil {
		s.logger.Error("Payment initialization failed",
			zap.Uint("order_id", newOrder.ID),
			zap.Error(err))

		// Rollback order creation if payment initialization fails
		if deleteErr := s.db.Delete(newOrder).Error; deleteErr != nil {
			s.logger.Error("Failed to rollback order", zap.Error(deleteErr))
		}

		return nil, fmt.Errorf("payment initialization failed: %w", err)
	}

	// Convert to DTO for response
	response := helpers.ToOrderResponse(newOrder)
	response.PaymentAuthorizationURL = paymentResp.AuthorizationURL
	response.PaymentReference = paymentResp.TransactionID

	// Send order confirmation email to customer
	if s.emailService != nil {
		// Prepare order items for email
		var emailItems []map[string]interface{}
		for _, item := range newOrder.OrderItems {
			emailItems = append(emailItems, map[string]interface{}{
				"Name":     item.Product.Name,
				"Quantity": item.Quantity,
				"Price":    fmt.Sprintf("₦%.2f", item.Price),
			})
		}

		// Send email in a goroutine to avoid blocking the response
		go func() {
			emailData := map[string]interface{}{
				"CustomerName":    user.Name,
				"OrderID":         fmt.Sprintf("%d", newOrder.ID),
				"OrderDate":       newOrder.CreatedAt.Format("January 2, 2006"),
				"TotalAmount":     fmt.Sprintf("₦%.2f", newOrder.TotalAmount.InexactFloat64()),
				"Items":           emailItems,
				"OrderDetailsURL": fmt.Sprintf("https://perthmarketplace.com/orders/%d", newOrder.ID),
				"MarketplaceURL":  "https://perthmarketplace.com",
			}

			if err := s.emailService.SendOrderConfirmation(user.Email, fmt.Sprintf("%d", newOrder.ID), emailData); err != nil {
				s.logger.Error("Failed to send order confirmation email", zap.Error(err))
			}
		}()

		// Send notification emails to merchants
		go func() {
			// Group items by merchant
			merchantItems := make(map[string][]map[string]interface{})
			merchantEmails := make(map[string]string)

			for _, item := range newOrder.OrderItems {
				merchantID := item.MerchantID
				if _, exists := merchantItems[merchantID]; !exists {
					merchantItems[merchantID] = []map[string]interface{}{}
					// Get merchant email (in a real implementation, you would fetch this from the database)
					merchantEmails[merchantID] = fmt.Sprintf("merchant-%s@perthmarketplace.com", merchantID)
				}

				merchantItems[merchantID] = append(merchantItems[merchantID], map[string]interface{}{
					"Name":     item.Product.Name,
					"Quantity": item.Quantity,
					"Price":    fmt.Sprintf("₦%.2f", item.Price),
				})
			}

			// Send email to each merchant
			for merchantID, items := range merchantItems {
				emailData := map[string]interface{}{
					"MerchantName":         fmt.Sprintf("Merchant %s", merchantID),
					"OrderID":              fmt.Sprintf("%d", newOrder.ID),
					"OrderDate":            newOrder.CreatedAt.Format("January 2, 2006"),
					"TotalAmount":          fmt.Sprintf("₦%.2f", newOrder.TotalAmount.InexactFloat64()),
					"Items":                items,
					"MerchantDashboardURL": "https://perthmarketplace.com/merchant/dashboard",
				}

				merchantEmail := merchantEmails[merchantID]
				if err := s.emailService.SendMerchantOrderNotification(merchantEmail, fmt.Sprintf("%d", newOrder.ID), emailData); err != nil {
					s.logger.Error("Failed to send merchant order notification email", zap.Error(err))
				}
			}
		}()
	}

	s.logger.Info("Order created successfully",
		zap.Uint("order_id", newOrder.ID),
		zap.Uint("user_id", userID),
		zap.String("payment_reference", paymentResp.TransactionID))

	return response, nil
}

// GetOrder retrieves a single order by its ID.
func (s *OrderService) GetOrder(ctx context.Context, id uint) (*models.Order, error) {
	if id == 0 {
		return nil, errors.New("invalid order ID")
	}
	// The repository already preloads necessary associations.
	return s.orderRepo.FindByID(ctx, id)
	//return s.orderRepo.FindByID(id)

}

// GetOrdersByUserID retrieves all orders for a user
// func (s *OrderService) GetOrdersByUserID(userID uint) ([]models.Order, error) {
// 	if userID == 0 {
// 		return nil, errors.New("invalid user ID")
// 	}
// 	return s.orderRepo.FindByUserID(userID)
// }

// GetOrdersByMerchantID retrieves orders containing a merchant's products
func (s *OrderService) GetOrdersByMerchantID(ctx context.Context, merchantID uint) ([]models.Order, error) {
	if merchantID == 0 {
		return nil, errors.New("invalid merchant ID")
	}
	return s.orderRepo.FindByMerchantID(ctx, merchantID)
}

// GetMerchantOrders retrieves orders containing a merchant's products using string merchant ID
func (s *OrderService) GetMerchantOrders(ctx context.Context, merchantID string) ([]models.Order, error) {
	if merchantID == "" {
		return nil, errors.New("invalid merchant ID")
	}
	// Convert string merchant ID to uint if needed
	// For now, we'll assume the repository can handle string merchant IDs
	// If not, we'll need to modify the repository method
	return s.orderRepo.FindByMerchantID(ctx, merchantID)
}

// UpdateOrderStatus updates the status of an order
func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID uint, status string) (*models.Order, error) {
	if orderID == 0 {
		return nil, errors.New("invalid order ID")
	}
	if err := models.OrderStatus(status).Valid(); err != nil {
		return nil, err
	}

	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	order.Status = models.OrderStatus(status)
	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	return s.orderRepo.FindByID(ctx, orderID)
}

// CancelOrder orchestrates cancellation (business logic here)
func (s *OrderService) CancelOrder(ctx context.Context, orderID uint, userID uint, reason string) error {
	logger := s.logger.With(zap.String("operation", "CancelOrder"), zap.Uint("order_id", orderID), zap.Uint("user_id", userID))

	// Fetch order (ownership checked in repo for efficiency)
	order, err := s.orderRepo.FindByIDWithPreloads(ctx, orderID)
	if err != nil {
		logger.Error("Failed to fetch order", zap.Error(err))
		return err
	}
	if order.UserID != userID {
		logger.Warn("Unauthorized cancellation attempt")
		return ErrUnauthorizedOrder
	}
	if order.Status != models.OrderStatusPending { // Adjust enum as per model
		logger.Warn("Invalid status for cancellation", zap.String("status", string(order.Status)))
		return ErrInvalidOrderStatus
	}

	// Transaction for atomicity
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update order status
		if err := s.orderRepo.UpdateStatus(ctx, orderID, models.OrderStatusCancelled); err != nil {
			return err
		}

		// Unreserve inventory for items (no VariantID, so use ProductID + MerchantID)
		items, err := s.orderItemRepo.FindOrderItemsByOrderID(ctx, orderID)
		if err != nil {
			return err
		}
		for _, item := range items {
			inventory, err := s.inventoryRepo.FindByProductAndMerchant(ctx, item.ProductID, item.MerchantID)
			if err != nil {
				return fmt.Errorf("inventory lookup failed for product %s: %w", item.ProductID, err)
			}
			// Unreserve: Add back to Quantity, subtract from ReservedQuantity
			inventory.Quantity += item.Quantity
			if inventory.ReservedQuantity >= item.Quantity {
				inventory.ReservedQuantity -= item.Quantity
			} else {
				inventory.ReservedQuantity = 0
			}
			if err := s.inventoryRepo.UpdateInventory(ctx, inventory.ID, item.Quantity); err != nil { // Assume repo method for update
				return fmt.Errorf("failed to update inventory %s: %w", inventory.ID, err)
			}
		}

		// Initiate refund if paid
		// if order.Payment != nil && order.Payment.Status == "success" {
		// 	if err := s.paymentService.InitiateRefund(ctx, orderID); err != nil {
		// 		logger.Error("Refund initiation failed", zap.Error(err))
		// 		return ErrRefundFailed
		// 	}
		// }

		return nil
	})
	if err != nil {
		logger.Error("Transaction failed", zap.Error(err))
		return err
	}

	// Notifications (outside tx, fire-and-forget)
	// if err := s.notificationSvc.NotifyUser(ctx, userID, "Order Cancelled", fmt.Sprintf("Order %d cancelled: %s", orderID, reason)); err != nil {
	// 	logger.Warn("User notification failed", zap.Error(err)) // Soft fail
	// }
	// // For multi-vendor: Notify per merchant (stub; loop over items if needed)
	// for _, item := range items { // From earlier fetch
	// 	if err := s.notificationSvc.NotifyMerchant(ctx, item.MerchantID, "Order Item Cancelled", fmt.Sprintf("Item for order %d cancelled", orderID)); err != nil {
	// 		logger.Warn("Merchant notification failed", zap.Error(err))
	// 	}
	// }

	logger.Info("Order cancelled successfully")
	return nil
}

// AcceptOrderItem allows a merchant to accept an order item
func (s *OrderService) AcceptOrderItem(ctx context.Context, orderItemID uint, merchantID string) error {
	// Fetch the order item
	orderItem, err := s.orderItemRepo.FindByIDWithContext(ctx, orderItemID)
	if err != nil {
		return fmt.Errorf("failed to find order item: %w", err)
	}

	// Verify the merchant owns this order item
	if orderItem.MerchantID != merchantID {
		return errors.New("unauthorized: merchant does not own this order item")
	}

	// Validate status transition using state machine logic
	if err := orderItem.ValidateStatusTransition(models.FulfillmentStatusConfirmed); err != nil {
		return fmt.Errorf("invalid status transition: %w", err)
	}

	// Update the fulfillment status to Confirmed
	orderItem.FulfillmentStatus = models.FulfillmentStatusConfirmed
	if err := s.orderItemRepo.Update(orderItem); err != nil {
		return fmt.Errorf("failed to update order item status: %w", err)
	}

	return nil
}

// DeclineOrderItem allows a merchant to decline an order item
func (s *OrderService) DeclineOrderItem(ctx context.Context, orderItemID uint, merchantID string) error {
	// Fetch the order item
	orderItem, err := s.orderItemRepo.FindByIDWithContext(ctx, orderItemID)
	if err != nil {
		return fmt.Errorf("failed to find order item: %w", err)
	}

	// Verify the merchant owns this order item
	if orderItem.MerchantID != merchantID {
		return errors.New("unauthorized: merchant does not own this order item")
	}

	if err := orderItem.ValidateStatusTransition(models.FulfillmentStatusDeclined); err != nil {
		return fmt.Errorf("invalid status transition: %w", err)
	}
	// Update the fulfillment status to Declined
	orderItem.FulfillmentStatus = models.FulfillmentStatusDeclined
	if err := s.orderItemRepo.Update(orderItem); err != nil {
		return fmt.Errorf("failed to update order item status: %w", err)
	}

	return nil
}

// UpdateOrderItemToSentToAronovaHub allows a merchant to update an order item to "SentToAronovaHub" status
func (s *OrderService) UpdateOrderItemToSentToAronovaHub(ctx context.Context, orderItemID uint, merchantID string) error {
	// Fetch the order item
	orderItem, err := s.orderItemRepo.FindByIDWithContext(ctx, orderItemID)
	if err != nil {
		return fmt.Errorf("failed to find order item: %w", err)
	}

	// Verify the merchant owns this order item
	if orderItem.MerchantID != merchantID {
		return errors.New("unauthorized: merchant does not own this order item")
	}

	if err := orderItem.ValidateStatusTransition(models.FulfillmentStatusSentToAronovaHub); err != nil {
		return fmt.Errorf("invalid status transition: %w", err)
	}

	// Update the fulfillment status to SentToAronovaHub
	orderItem.FulfillmentStatus = models.FulfillmentStatusSentToAronovaHub
	if err := s.orderItemRepo.Update(orderItem); err != nil {
		return fmt.Errorf("failed to update order item status: %w", err)
	}

	return nil
}

func (s *OrderService) GetUserOrders(ctx context.Context, userID uint) ([]dto.OrdersResponse, error) {
	orders, err := s.orderRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var orderDTOs []dto.OrdersResponse
	for _, order := range orders {
		orderDTO := dto.OrdersResponse{
			ID:        order.ID,
			UserID:    order.UserID,
			Status:    dto.OrderStatus(order.Status),
			CreatedAt: order.CreatedAt,
			UpdatedAt: order.UpdatedAt,
		}
		for _, item := range order.OrderItems {
			itemDTO := dto.OrdersItemResponse{
				ID:      item.ID,
				OrderID: item.OrderID,
				//ProductID: item.ProductID,
				Quantity: uint(item.Quantity),
				//MerchantID: item.MerchantID,
			}
			if item.Product.ID != "" {
				itemDTO.Product = dto.OrderProductResponse{
					ID:          item.Product.ID,
					Name:        item.Product.Name,
					Description: item.Product.Description,
					Price:       item.Product.BasePrice.InexactFloat64(),
				}
				if len(item.Product.Media) > 0 {
					itemDTO.Product.Image = item.Product.Media[0].URL // Assume Media has URL field
				}
			}
			if item.Merchant.MerchantID != "" {
				itemDTO.Merchant = dto.OrderMerchantResponse{
					ID:        item.Merchant.ID,
					StoreName: item.Merchant.StoreName,
				}
			}
			orderDTO.OrderItems = append(orderDTO.OrderItems, itemDTO)
		}
		orderDTOs = append(orderDTOs, orderDTO)
	}
	return orderDTOs, nil
}
