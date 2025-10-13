package order

import (
	"context"
	"errors"
	"fmt"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"

	//"go.uber.org/zap"
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// OrderService provides business logic for handling orders.
type OrderService struct {
	orderRepo     *repositories.OrderRepository
	orderItemRepo *repositories.OrderItemRepository
	cartRepo      *repositories.CartRepository
	cartItemRepo  *repositories.CartItemRepository
	productRepo   *repositories.ProductRepository
	inventoryRepo *repositories.InventoryRepository
	logger      *zap.Logger
	validator   *validator.Validate
	db            *gorm.DB
}

// NewOrderService creates a new instance of OrderService.
func NewOrderService(
	orderRepo *repositories.OrderRepository,
	orderItemRepo *repositories.OrderItemRepository,
	cartRepo *repositories.CartRepository,
	cartItemRepo *repositories.CartItemRepository,
	productRepo *repositories.ProductRepository,
	inventoryRepo *repositories.InventoryRepository,
) *OrderService {
	return &OrderService{
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
		cartRepo:      cartRepo,
		cartItemRepo:  cartItemRepo,
		productRepo:   productRepo,
		inventoryRepo: inventoryRepo,
		db:            db.DB,
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



func (s *OrderService) CreateOrder(ctx context.Context, userID uint) (*dto.OrderResponse, error) {
    if userID == 0 {
        return nil, errors.New("invalid user ID")
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
    var totalAmount float64

    // Use a transaction to ensure atomicity
    err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // Calculate total and create order items
        var orderItems []models.OrderItem
        for _, item := range cart.CartItems {
            price := item.Product.FinalPrice.InexactFloat64()
            if item.VariantID != nil && item.Variant != nil {
                price = item.Variant.FinalPrice.InexactFloat64() // Use variant price if available
            }
            totalAmount += float64(item.Quantity) * price
            orderItems = append(orderItems, models.OrderItem{
                ProductID:         item.ProductID,
                MerchantID:        item.Product.MerchantID,
                Quantity:          item.Quantity,
                Price:             price,
                FulfillmentStatus: models.FulfillmentStatusNew,
                // ID is omitted to let GORM/database auto-generate
            })

            // Update inventory: move reserved stock to committed stock
            inventoryQuery := tx.Model(&models.Inventory{}).Where("product_id = ?", item.ProductID)
            if item.VariantID != nil {
                inventoryQuery = inventoryQuery.Where("variant_id = ?", *item.VariantID)
            }
            if err := inventoryQuery.Updates(map[string]interface{}{
                "quantity":          gorm.Expr("quantity - ?", item.Quantity),
                "reserved_quantity": gorm.Expr("reserved_quantity - ?", item.Quantity),
            }).Error; err != nil {
                return fmt.Errorf("failed to commit stock for product %s: %w", item.ProductID, err)
            }
        }

        // Create the order
        newOrder = &models.Order{
            UserID:      userID,
            TotalAmount: decimal.NewFromFloat(totalAmount),
            Status:      models.OrderStatusPending,
            OrderItems:  orderItems,
        }
        if err := tx.Create(newOrder).Error; err != nil {
            return fmt.Errorf("failed to create order: %w", err)
        }

        // Associate order items with the new order ID
        for i := range orderItems {
            orderItems[i].OrderID = newOrder.ID
            orderItems[i].ID = 0 // Explicitly reset ID to ensure auto-generation
        }
        if err := tx.Create(&orderItems).Error; err != nil {
            return fmt.Errorf("failed to create order items: %w", err)
        }

        // Manually associate for the response DTO
        newOrder.OrderItems = orderItems

        // Clear cart items
        if err := tx.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
            return fmt.Errorf("failed to clear cart items: %w", err)
        }

        // Mark cart as converted
        cart.Status = models.CartStatusConverted
        if err := tx.Save(cart).Error; err != nil {
            return fmt.Errorf("failed to update cart status: %w", err)
        }

        return nil
    })

    if err != nil {
        //s.logger.Error("Transaction failed", zap.Error(err))
        return nil, fmt.Errorf("transaction failed: %w", err)
    }

    // Convert to DTO for response
    orderResponse := &dto.OrderResponse{
        ID:         newOrder.ID,
        UserID:     newOrder.UserID,
        Status:     dto.OrderStatus(newOrder.Status),
        OrderItems: make([]dto.OrderItemResponse, len(newOrder.OrderItems)),
		CreatedAt: newOrder.CreatedAt,
		UpdatedAt: newOrder.UpdatedAt,
    }
    for i, item := range newOrder.OrderItems {
        orderResponse.OrderItems[i] = dto.OrderItemResponse{
            ProductID: fmt.Sprint(item.ProductID),
			Name: item.Product.Name,
            Quantity:  item.Quantity,
            Price:     item.Price,
        }
    }
	//  user, _ := userRepo.FindByID(userID)  // Assume userRepo injected
    // authURL, ref, err := paymentService.InitiateTransaction(ctx, order, user.Email)
    // if err != nil {
    //     return nil, err
    // }
    // Return order with authURL in response

    return orderResponse, nil
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
func (s *OrderService) GetOrdersByMerchantID(ctx context.Context ,merchantID uint) ([]models.Order, error) {
	if merchantID == 0 {
		return nil, errors.New("invalid merchant ID")
	}
	return s.orderRepo.FindByMerchantID(ctx,merchantID)
}

// UpdateOrderStatus updates the status of an order
func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID uint, status string) (*models.Order, error) {
	if orderID == 0 {
		return nil, errors.New("invalid order ID")
	}
	if err := models.OrderStatus(status).Valid(); err != nil {
		return nil, err
	}

	order, err := s.orderRepo.FindByID(ctx,orderID)
	if err != nil {
		return nil, err
	}

	order.Status = models.OrderStatus(status)
	if err := s.orderRepo.Update(ctx ,order); err != nil {
		return nil, err
	}

	return s.orderRepo.FindByID(ctx ,orderID)
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
				ID:        item.ID,
				OrderID:   item.OrderID,
				ProductID: item.ProductID,
				Quantity:  uint(item.Quantity),
				MerchantID: item.MerchantID,
			}
			if item.Product.ID != "" {
				itemDTO.Product = dto.OrderProductResponse{
					ID:          item.Product.ID,
					Name:        item.Product.Name,
					Description: item.Product.Description,
					Price:       item.Product.BasePrice.InexactFloat64(),
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

