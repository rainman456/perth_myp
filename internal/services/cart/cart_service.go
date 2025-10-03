package cart

import (
	//"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/api/dto" // Assuming dto.BulkUpdateRequest is defined here with ProductID string, Quantity int
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"
	"context"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
	//"gorm.io/gorm/logger"
)

var (
	ErrInvalidUserID     = errors.New("invalid user ID")
	ErrInvalidQuantity   = errors.New("quantity must be positive")
	ErrProductNotFound   = errors.New("product not found")
	ErrInventoryNotFound = errors.New("inventory not found")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrTransactionFailed   = errors.New("transaction failed")
)

type CartService struct {
	cartRepo      *repositories.CartRepository
	cartItemRepo  *repositories.CartItemRepository
	productRepo   *repositories.ProductRepository
	inventoryRepo *repositories.InventoryRepository
	logger        *zap.Logger
	validator     *validator.Validate
}

func NewCartService(cartRepo *repositories.CartRepository, cartItemRepo *repositories.CartItemRepository, productRepo *repositories.ProductRepository, inventoryRepo *repositories.InventoryRepository, logger *zap.Logger) *CartService {
	return &CartService{
		cartRepo:      cartRepo,
		cartItemRepo:  cartItemRepo,
		productRepo:   productRepo,
		inventoryRepo: inventoryRepo,
		logger:        logger,
		validator:     validator.New(),
	}
}

// GetActiveCart retrieves or creates an active cart for a user
func (s *CartService) GetActiveCart(ctx context.Context, userID uint) (*dto.CartResponse, error) {
	if userID == 0 {
		return nil, ErrInvalidUserID
	}
	cart, err := s.cartRepo.FindActiveCart(ctx, userID)
	// if err == nil {
	// 	return cart, nil
	// }
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error("Failed to find active cart", zap.Uint("user_id", userID), zap.Error(err))
		return nil, fmt.Errorf("db error: %w", err)
	}

	cart = &models.Cart{UserID: userID, Status: models.CartStatusActive}
	if err := s.cartRepo.Create(ctx, cart); err != nil {
		s.logger.Error("Failed to create cart", zap.Error(err))
		return nil, fmt.Errorf("create failed: %w", err)
	}
	//return s.cartRepo.FindByID(ctx, cart.ID)
	cart, err = s.cartRepo.FindByID(ctx, cart.ID)
	if err != nil {
		s.logger.Error("Failed to get active cart", zap.Error(err))
		return nil, fmt.Errorf("failed to get active cart: %w", err)
	}
	response := &dto.CartResponse{
	ID:        cart.ID,
	UserID:    cart.UserID,
	Status:    cart.Status,
	Items:     make([]dto.CartItemResponse, len(cart.CartItems)),
	Total:     cart.GrandTotal, // Assuming decimal.Decimal
	CreatedAt: cart.CreatedAt,
	UpdatedAt: cart.UpdatedAt,
}
for i, item := range cart.CartItems {
	response.Items[i] = dto.CartItemResponse{
			ID:        item.ID,
		ProductID: item.ProductID,
		VariantID: item.VariantID, // Fixed from m.URL
		Quantity:  item.Quantity,
		Subtotal:  item.Cart.SubTotal,
	}
}
return response, nil

}

// func (s *CartService) GetCart(ctx context.Context, userID uint) (*models.Cart, error) {
// 	cart, err := s.GetActiveCart(ctx, userID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// Fixed: Use model method
// 	cart.ComputeTotals()
// 	s.logger.Info("Cart fetched", zap.Uint("user_id", userID), zap.Float64("total", cart.GrandTotal))
// 	return cart, nil
// }

// AddItemToCart adds a product to the user's active cart
/*
func (s *CartService) AddItemToCart(ctx context.Context, userID uint, quantity uint, productID string) (*models.Cart, error) {
	if userID == 0 {
		return nil, ErrInvalidUserID
	}
	if productID == "" {
		return nil, errors.New("invalid product ID")
	}
	if quantity == 0 {
		return nil, ErrInvalidQuantity
	}

	cart, err := s.GetActiveCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return nil, ErrProductNotFound
	}

	inventory, err := s.inventoryRepo.FindByProductID(ctx,productID,product.MerchantID)
	if err != nil {
		return nil, ErrInventoryNotFound
	}
	if inventory.StockQuantity < int(quantity) {
		return nil, ErrInsufficientStock
	}

	err = db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		cartItem, err := s.cartItemRepo.FindByProductIDAndCartID(ctx, productID, cart.ID)
		newQty := quantity
		if err == nil {
			newQty += uint(cartItem.Quantity)
			if inventory.StockQuantity < int(newQty) {
				return ErrInsufficientStock
			}
			if err := s.cartItemRepo.UpdateQuantityWithReservation(ctx, cartItem.ID, int(newQty), inventory.ID); err != nil {
				return err
			}
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			cartItem = &models.CartItem{
				CartID:     cart.ID,
				ProductID:  productID,
				Quantity:   int(quantity),
				MerchantID: product.MerchantID,
			}
			if err := s.cartItemRepo.Create(ctx, cartItem); err != nil {
				return err
			}
			if err := s.inventoryRepo.UpdateInventoryQuantity(ctx, inventory.ID, -int(quantity)); err != nil {
				return err
			}
		} else {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	// Reload
	return s.cartRepo.FindByID(ctx, cart.ID)
}
*/


/*
func (s *CartService) AddItemToCart(ctx context.Context, userID uint, quantity uint, productID string) (*models.Cart, error) {
	if userID == 0 {
		return nil, ErrInvalidUserID
	}
	if quantity == 0 {
		return nil, ErrInvalidQuantity
	}

	cart, err := s.GetActiveCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	product, err := s.productRepo.FindByID(ctx,productID)
	if err != nil {
		return nil, err
	}

	inventory, err := s.inventoryRepo.FindByProductID(ctx, productID, product.MerchantID)
	if err != nil {
		return nil, ErrInventoryNotFound
	}

	if inventory.Quantity < int(quantity) {
		return nil, ErrInsufficientStock
	}

	// Assuming no variant for simplicity; pass nil for variantID
	existing, err := s.cartItemRepo.FindByProductIDAndCartID(ctx, productID, nil, cart.ID)
	if err == nil {
		// Existing item: increment quantity
		newQty := existing.Quantity + int(quantity)
		if newQty > inventory.Quantity {
			return nil, ErrInsufficientStock
		}
		err = s.cartItemRepo.UpdateQuantityWithReservation(ctx, existing.ID, newQty, inventory.ID)
		if err != nil {
			return nil, err
		}
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// New item
		cartItem := &models.CartItem{
			CartID:    cart.ID,
			ProductID: productID,
			VariantID: nil,
			Quantity:  int(quantity),
			//PriceSnapshot: product.BasePrice.InexactFloat64(),
			MerchantID: product.MerchantID,
		}
		if err := s.cartItemRepo.Create(ctx, cartItem); err != nil {
			return nil, err
		}
		if err := s.inventoryRepo.UpdateInventoryQuantity(ctx, inventory.ID, -int(quantity)); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	return s.cartRepo.FindByID(ctx, cart.ID)
}
*/




/*
func (s *CartService) AddItemToCart(ctx context.Context, userID uint, quantity uint, productID string, variantID *string) (*models.Cart, error) {
    if userID == 0 {
        return nil, ErrInvalidUserID
    }
    if productID == "" {
        return nil, errors.New("invalid product ID")
    }
    if quantity == 0 {
        return nil, ErrInvalidQuantity
    }

    cart, err := s.GetActiveCart(ctx, userID)
    if err != nil {
        return nil, err
    }

    product, err := s.productRepo.FindByID(ctx, productID)
    if err != nil {
        return nil, ErrProductNotFound
    }

    var inventory *models.VendorInventory
    var priceSnapshot decimal.Decimal = product.BasePrice  // Default to base

    if variantID != nil {
        // Variant product
        variant, err := s.variantRepo.FindByID(ctx, *variantID)
        if err != nil {
            return nil, ErrInvalidVariant
        }
        if variant.ProductID != product.ID {
            return nil, errors.New("variant does not belong to product")
        }
        inventory, err = s.inventoryRepo.FindByVariantID(ctx, *variantID, product.MerchantID)
        if err != nil {
            return nil, ErrInventoryNotFound
        }
        priceSnapshot = product.BasePrice.Add(variant.PriceAdjustment)
    } else {
        // Simple product
        if len(product.Variants) > 0 {
            return nil, errors.New("variant required for this product")
        }
        inventory, err = s.inventoryRepo.FindByProductID(ctx, productID, product.MerchantID)
        if err != nil {
            return nil, ErrInventoryNotFound
        }
    }

    // Quick pre-check
    available := inventory.Quantity - inventory.ReservedQuantity
    if available < int(quantity) {
        return nil, ErrInsufficientStock
    }

    // Transaction
    err = db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        var freshInv *models.VendorInventory
        var freshErr error

        if variantID != nil {
            // Re-fetch variant inventory inside tx for freshness
            freshInv, freshErr = s.inventoryRepo.FindByVariantID(tx, *variantID, product.MerchantID)
        } else {
            // Re-fetch product inventory inside tx for freshness
            freshInv, freshErr = s.inventoryRepo.FindByProductID(tx, productID, product.MerchantID)
        }

        if freshErr != nil {
            return freshErr
        }
        available = freshInv.Quantity - freshInv.ReservedQuantity
        if available < int(quantity) {
            return ErrInsufficientStock
        }

        // Existing item check (use tx)
        cartItem, err := s.cartItemRepo.FindByProductAndVariant(tx, productID, variantID, cart.ID)
        newQty := int(quantity)

        if err == nil {
            newQty += cartItem.Quantity
            if available < newQty {
                return ErrInsufficientStock
            }
            // Update with price if needed (for MVP, assume snapshot set on create)
            return s.cartItemRepo.UpdateQuantityWithReservation(tx, cartItem.ID, newQty, freshInv.ID)
        } else if errors.Is(err, gorm.ErrRecordNotFound) {
            cartItem = &models.CartItem{
                CartID:        cart.ID,
                ProductID:     productID,
                VariantID:     variantID,  // Nil-safe
                Quantity:      newQty,
                PriceSnapshot: priceSnapshot,  // Set once
                MerchantID:    product.MerchantID,
            }
            if err := s.cartItemRepo.Create(tx, cartItem); err != nil {
                return err
            }
            // Reserve (for create, delta = quantity)
            return s.inventoryRepo.ReserveStock(tx, freshInv.ID, newQty)
        }
        return err
    })
    if err != nil {
        s.logger.Error("Add to cart failed", zap.Error(err), zap.Uint("user_id", userID))
        return nil, err
    }

    return s.cartRepo.FindByID(ctx, cart.ID)
}


*/




















func (s *CartService) AddItemToCart(ctx context.Context, userID uint, quantity int, productID string, variantID *string) (*dto.CartResponse, error) {
	if userID == 0 {
		return nil, ErrInvalidUserID
	}
	if quantity <= 0 {
		return nil, ErrInvalidQuantity
	}

	cart, err := s.GetActiveCart(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get active cart", zap.Uint("user_id", userID), zap.Error(err))
		return nil, err
	}

	// Fetch product with preloaded Variants.Inventory and SimpleInventory
	product, err := s.productRepo.FindByID(ctx, productID, "Variants.Inventory", "SimpleInventory")
	if err != nil {
		s.logger.Error("Product not found", zap.String("product_id", productID), zap.Error(err))
		return nil, ErrProductNotFound
	}
	if product.DeletedAt.Valid {
		s.logger.Error("Product is soft-deleted", zap.String("product_id", productID))
		return nil, ErrProductNotFound
	}

	// Determine inventory: focus on variants if they exist, else simple
	var inventory *models.Inventory
	var price decimal.Decimal = product.BasePrice
	//var varID string
	if variantID != nil && len(product.Variants) > 0 {
		//varID = *variantID
		for _, v := range product.Variants {
			if v.ID == *variantID && v.IsActive {
				inventory = &v.Inventory
				price = price.Add(v.PriceAdjustment)
				break
			}
		}
	} else if variantID == nil && product.SimpleInventory != nil {
		inventory = product.SimpleInventory
	} else {
		s.logger.Error("Inventory not found", zap.String("product_id", productID), zap.Stringp("variant_id", variantID))
		return nil, ErrInventoryNotFound
	}
	if inventory == nil {
		s.logger.Error("No valid inventory", zap.String("product_id", productID), zap.Stringp("variant_id", variantID))
		return nil, ErrInventoryNotFound
	}

	// Check available stock
	available := inventory.Quantity - inventory.ReservedQuantity
	if available < quantity {
		s.logger.Warn("Insufficient stock", zap.Int("available", available), zap.Int("requested", quantity))
		return nil, ErrInsufficientStock
	}

	// Transaction: Update cart item and reserve inventory
	err = db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Find existing cart item
		existing, err := s.cartItemRepo.FindByProductIDAndCartID(ctx, productID, nil,cart.ID)
		if err == nil {
			// Update existing item
			newQty := existing.Quantity + quantity
			if newQty > available {
				return ErrInsufficientStock
			}
			if err := s.cartItemRepo.UpdateQuantityWithReservation(ctx, existing.ID, newQty, inventory.ID); err != nil {
				return fmt.Errorf("failed to update cart item: %w", err)
			}
			return nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("failed to check existing cart item: %w", err)
		}

		// Create new cart item
		cartItem := &models.CartItem{
			CartID:    cart.ID,
			ProductID: productID,
			VariantID: variantID, // Assume VariantID is *string in model
			Quantity:  quantity,
			MerchantID: product.MerchantID,
		}
		if err := s.cartItemRepo.Create(ctx, cartItem); err != nil {
			return fmt.Errorf("failed to create cart item: %w", err)
		}
		inventory.ReservedQuantity += quantity // Manual update since method undefined
		if err := tx.Save(inventory).Error; err != nil {
			return fmt.Errorf("failed to reserve inventory: %w", err)
		}
		return nil
	})
	if err != nil {
		s.logger.Error("Transaction failed", zap.Error(err))
		return nil, fmt.Errorf("%w: %v", ErrTransactionFailed, err)
	}

	// Return updated cart
	updatedCart, err := s.cartRepo.FindByID(ctx, cart.ID)
    if err != nil {
        s.logger.Error("Failed to fetch updated cart", zap.Uint("cart_id", cart.ID), zap.Error(err))
        return nil, err
    }
    // Fix: Preload CartItems with related data
    if err := db.DB.WithContext(ctx).
        Preload("CartItems.Product.Media").
        Preload("CartItems.Product.Variants.Inventory").
        Preload("CartItems.Variant").
        Find(updatedCart).Error; err != nil {
        s.logger.Error("Failed to preload cart items", zap.Error(err))
        return nil, err
    }
   // return updatedCart, nil
	response := &dto.CartResponse{
		ID: updatedCart.ID,
		UserID: updatedCart.UserID,
		Status: updatedCart.Status,
		Items: make([]dto.CartItemResponse, len(updatedCart.CartItems)),
		Total: updatedCart.GrandTotal,
		CreatedAt: updatedCart.CreatedAt,
		UpdatedAt: updatedCart.UpdatedAt,
	}
	for i, item := range updatedCart.CartItems {
	response.Items[i] = dto.CartItemResponse{
		ID:        item.ID,
		ProductID: item.ProductID,
		VariantID: item.VariantID, // Fixed from m.URL
		Quantity:  item.Quantity,
		Subtotal:  item.Cart.SubTotal,

	}
}
return response, nil
}











// UpdateCartItemQuantity updates the quantity of a cart item
func (s *CartService) UpdateCartItemQuantity(ctx context.Context, cartItemID uint, quantity int) (*models.Cart, error) {
	if cartItemID == 0 {
		return nil, errors.New("invalid cart item ID")
	}
	if quantity <= 0 {
		return nil, ErrInvalidQuantity
	}

	// load cart item (contains MerchantID and ProductID)
	cartItem, err := s.cartItemRepo.FindByID(ctx, cartItemID)
	if err != nil {
		return nil, repositories.ErrCartItemNotFound
	}

	// ensure we have a merchantID to scope the inventory lookup
	merchantID := cartItem.MerchantID
	if merchantID == "" {
		// fallback: fetch product to get merchant (shouldn't usually happen if cart items store merchant)
		prod, perr := s.productRepo.FindByID(ctx,cartItem.ProductID)
		if perr != nil {
			return nil, ErrInventoryNotFound
		}
		merchantID = prod.MerchantID
	}

	// NOTE: FindByProductID signature is (ctx, productID, merchantID)
	inventory, err := s.inventoryRepo.FindByProductID(ctx, cartItem.ProductID, merchantID)
	if err != nil {
		return nil, ErrInventoryNotFound
	}

	// model field is Quantity (not StockQuantity)
	if inventory.Quantity < quantity {
		return nil, ErrInsufficientStock
	}

	// UpdateQuantityWithReservation now expects vendor inventory ID as string
	if err := s.cartItemRepo.UpdateQuantityWithReservation(ctx, cartItemID, quantity, inventory.ID); err != nil {
		return nil, err
	}

	return s.cartRepo.FindByID(ctx, cartItem.CartID)
}

func (s *CartService) RemoveCartItem(ctx context.Context, cartItemID uint) (*models.Cart, error) {
	if cartItemID == 0 {
		return nil, errors.New("invalid cart item ID")
	}

	cartItem, err := s.cartItemRepo.FindByID(ctx, cartItemID)
	if err != nil {
		return nil, repositories.ErrCartItemNotFound
	}

	// use the merchant stored on the cart item to find the correct vendor inventory
	merchantID := cartItem.MerchantID
	if merchantID == "" {
		// fallback: fetch product to get merchant
		prod, perr := s.productRepo.FindByID(ctx,cartItem.ProductID)
		if perr != nil {
			return nil, ErrInventoryNotFound
		}
		merchantID = prod.MerchantID
	}

	// pass ctx and merchantID as required by repo
	inventory, err := s.inventoryRepo.FindByProductID(ctx, cartItem.ProductID, merchantID)
	if err != nil {
		return nil, ErrInventoryNotFound
	}

	// DeleteWithUnreserve expects vendor inventory ID as string
	if err := s.cartItemRepo.DeleteWithUnreserve(ctx, cartItemID, inventory.ID); err != nil {
		return nil, err
	}

	return s.cartRepo.FindByID(ctx, cartItem.CartID)
}





func (s *CartService) GetCartItemByID(ctx context.Context, cartItemID uint) (*models.CartItem, error) {
	if cartItemID == 0 {
		return nil, errors.New("invalid cart item ID")
	}
	return s.cartItemRepo.FindByID(ctx, cartItemID)
}





// ClearCart, BulkAddItems ... (add ctx to all calls; stub Bulk if not used)
func (s *CartService) ClearCart(ctx context.Context, userID uint) error {
	cart, err := s.cartRepo.FindActiveCart(ctx, userID)
	if err != nil {
		return err
	}
	items, err := s.cartItemRepo.FindByCartID(ctx, cart.ID)
	if err != nil {
		return err
	}
	for _, item := range items {
		s.RemoveCartItem(ctx, item.ID)
	}
	cart.Status = models.CartStatusAbandoned
	return s.cartRepo.Update(ctx, cart)
}

// BulkAddItems stub (implement as needed; fixed DTO)
// func (s *CartService) BulkAddItems(ctx context.Context, userID uint, items []dto.BulkUpdateRequest) (*models.Cart, error) {
// 	// Validation loop...
// 	for _, item := range items {
// 		if err := s.validator.Struct(&item); err != nil {
// 			return nil, err
// 		}
// 		// Add each (loop AddItemToCart)
// 		_, err := s.AddItemToCart(ctx, userID, uint(item.Quantity), item.ProductID)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	return s.GetCart(ctx, userID)
// }












// func (s *CartService) BulkAddItems(ctx context.Context, userID uint, items dto.BulkUpdateRequest) (*models.Cart, error) {
// 	if userID == 0 {
// 		return nil, ErrInvalidUserID
// 	}
// 	if err := s.validator.Struct(&items); err != nil {
// 		return nil, fmt.Errorf("validation failed: %w", err)
// 	}

// 	// cart, err := s.GetActiveCart(ctx, userID)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	for _, item := range items.Items {
// 		// Convert uint ProductID to string for consistency
// 		productID := fmt.Sprint(item.ProductID)
// 		if _, err := s.AddItemToCart(ctx, userID, uint(item.Quantity), productID); err != nil {
// 			return nil, fmt.Errorf("failed to add item %s: %w", productID, err)
// 		}
// 	}
// 	return s.GetCart(ctx, userID)
// }






func (s *CartService) BulkAddItems(ctx context.Context, userID uint, items dto.BulkUpdateRequest) (*models.Cart, error) {
	if userID == 0 {
		return nil, ErrInvalidUserID
	}
	if len(items.Items) == 0 {
		return nil, errors.New("no items provided")
	}
	if err := s.validator.Struct(&items); err != nil {
		s.logger.Error("Validation failed", zap.Uint("user_id", userID), zap.Error(err))
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	cart, err := s.GetActiveCart(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get active cart", zap.Uint("user_id", userID), zap.Error(err))
		return nil, err
	}

	err = db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i, item := range items.Items {
			if _, err := s.AddItemToCart(ctx, userID, item.Quantity, item.ProductID, item.VariantID); err != nil {
				s.logger.Error("Failed to add item", zap.String("product_id", item.ProductID), zap.Stringp("variant_id", item.VariantID), zap.Error(err))
				return fmt.Errorf("failed to add item %d (product %s): %w", i+1, item.ProductID, err)
			}
		}
		return nil
	})
	if err != nil {
		s.logger.Error("Transaction failed", zap.Uint("user_id", userID), zap.Error(err))
		return nil, fmt.Errorf("%w: %v", ErrTransactionFailed, err)
	}

	updatedCart, err := s.cartRepo.FindByID(ctx, cart.ID)
	if err != nil {
		s.logger.Error("Failed to fetch updated cart", zap.Uint("cart_id", cart.ID), zap.Error(err))
		return nil, err
	}
	return updatedCart, nil
}