

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
	"gorm.io/gorm"
)

// OrderService provides business logic for handling orders.
type OrderService struct {
	orderRepo     *repositories.OrderRepository
	orderItemRepo *repositories.OrderItemRepository
	cartRepo      *repositories.CartRepository
	cartItemRepo  *repositories.CartItemRepository
	productRepo   *repositories.ProductRepository
	db            *gorm.DB
}

// NewOrderService creates a new instance of OrderService.
func NewOrderService(
	orderRepo *repositories.OrderRepository,
	orderItemRepo *repositories.OrderItemRepository,
	cartRepo *repositories.CartRepository,
	cartItemRepo *repositories.CartItemRepository,
	productRepo *repositories.ProductRepository,
) *OrderService {
	return &OrderService{
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
		cartRepo:      cartRepo,
		cartItemRepo:  cartItemRepo,
		productRepo:   productRepo,
		db:            db.DB,
	}
}

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
/*
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
			price := item.Product.BasePrice.InexactFloat64()
			totalAmount += float64(item.Quantity) * price
			orderItems = append(orderItems, models.OrderItem{
				ProductID:  item.ProductID,
				MerchantID: item.Product.MerchantID,
				Quantity:   item.Quantity,
				Price:      price,
				FulfillmentStatus: models.FulfillmentStatusNew,
			})

			// Here you would typically move reserved stock to committed stock.
			// For now, we assume cart reservation handled this.
			// We'll just update the main inventory.
			// This logic might need to be more robust depending on inventory strategy.
			if err := tx.Model(&models.Inventory{}).
				Where("product_id = ?", item.ProductID).
				Updates(map[string]interface{}{
					"quantity":          gorm.Expr("quantity - ?", item.Quantity),
					"reserved_quantity": gorm.Expr("reserved_quantity - ?", item.Quantity),
				}).Error; err != nil {
				return fmt.Errorf("failed to commit stock for product %s: %w", item.ProductID, err)
			}
		}

		// Create the order
		newOrder = &models.Order{
			UserID:      userID,
			TotalAmount: totalAmount,
			Status:      models.OrderStatusPending,
			OrderItems:  orderItems,
		}
		if err := tx.Create(newOrder).Error; err != nil {
			return err
		}

		// Associate order items with the new order ID and create them
		for i := range orderItems {
			orderItems[i].OrderID = newOrder.ID
		}
		if err := tx.Create(&orderItems).Error; err != nil {
			return err
		}

		// Manually associate for the response DTO
		newOrder.OrderItems = orderItems

		// Clear cart items
		if err := tx.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
			return err
		}

		// Mark cart as converted
		cart.Status = models.CartStatusConverted
		if err := tx.Save(cart).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Convert to DTO for response
	orderResponse := &dto.OrderResponse{
		ID:         newOrder.ID,
		UserID:     newOrder.UserID,
		Status:     string(newOrder.Status),
		OrderItems: make([]dto.OrderItemResponse, len(newOrder.OrderItems)),
	}
	for i, item := range newOrder.OrderItems {
		orderResponse.OrderItems[i] = dto.OrderItemResponse{
			ProductID: fmt.Sprint(item.ProductID),
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	return orderResponse, nil
}
*/


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
            price := item.Product.BasePrice.InexactFloat64()
            if item.VariantID != nil && item.Variant != nil {
                price = item.Variant.TotalPrice.InexactFloat64() // Use variant price if available
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
            TotalAmount: totalAmount,
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
        Status:     string(newOrder.Status),
        OrderItems: make([]dto.OrderItemResponse, len(newOrder.OrderItems)),
    }
    for i, item := range newOrder.OrderItems {
        orderResponse.OrderItems[i] = dto.OrderItemResponse{
            ProductID: fmt.Sprint(item.ProductID),
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
