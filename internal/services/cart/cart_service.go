package cart

import (
	//"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/api/dto" // Assuming dto.BulkUpdateRequest is defined here with ProductID string, Quantity int
	"api-customer-merchant/internal/api/helpers"
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
	"gorm.io/gorm/clause"
	//"gorm.io/gorm/logger"
)

var (
	ErrInvalidUserID     = errors.New("invalid user ID")
	ErrInvalidQuantity   = errors.New("quantity must be positive")
	ErrProductNotFound   = errors.New("product not found")
	ErrInventoryNotFound = errors.New("inventory not found")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrTransactionFailed = errors.New("transaction failed")
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
	// Error only on unexpected DB issues (not "not found")
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error("Failed to query active cart", zap.Uint("user_id", userID), zap.Error(err))
		return nil, fmt.Errorf("db error: %w", err)
	}

	// If not found (or nil), create new
	if errors.Is(err, gorm.ErrRecordNotFound) || cart == nil {
		newCart := &models.Cart{UserID: userID, Status: models.CartStatusActive}
		if createErr := s.cartRepo.Create(ctx, newCart); createErr != nil {
			s.logger.Error("Failed to create cart", zap.Error(createErr))
			return nil, fmt.Errorf("create failed: %w", createErr)
		}
		// Fetch created (with ID now set)
		cart, err = s.cartRepo.FindByID(ctx, newCart.ID)
		if err != nil || cart == nil {
			s.logger.Error("Failed to fetch created cart", zap.Error(err))
			return nil, fmt.Errorf("failed to get active cart: %w", err)
		}
		s.logger.Info("Created new active cart", zap.Uint("cart_id", cart.ID))
	}
// 	response := &dto.CartResponse{
//     ID:        cart.ID,
//     UserID:    cart.UserID,
//     Status:    cart.Status,
//     Items:     make([]dto.CartItemResponse, len(cart.CartItems)),
//     Total:     0,  // Will compute sum of subtotals
//     CreatedAt: cart.CreatedAt,
//     UpdatedAt: cart.UpdatedAt,
// }
// for i, item := range cart.CartItems {
//     subtotal := 0.0
//     // Check if Product is preloaded (avoid empty struct issues)
//     if item.Product.ID != "" {  // Use ID as non-zero check (struct-safe)
//         price := item.Product.FinalPrice.InexactFloat64()  // Convert decimal.Decimal
//         if item.VariantID != nil && item.Variant != nil && item.Variant.ID != "" {
//             price += item.Variant.FinalPrice.InexactFloat64()  // Add adjustment
//         }
//         subtotal = float64(item.Quantity) * price
//     }
//     response.Items[i] = dto.CartItemResponse{
//         ID:        item.ID,
//         ProductID: item.ProductID,
// 		Name: item.Product.Name,
//         VariantID: item.VariantID,
//         Quantity:  item.Quantity,
//         Subtotal:  subtotal,  // Computed: quantity * (base + adjustment)
//     }
//     response.Total += subtotal  // Accumulate grand total
// }
// return response, nil
response := helpers.ToCartResponse(cart)
	return response, nil

}

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

    // Determine inventory and price
    var inventory *models.Inventory
    var price decimal.Decimal = product.BasePrice
    if variantID != nil && len(product.Variants) > 0 {
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

    // Transaction: Update cart item and reserve inventory
    err = db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // Lock inventory
        var lockedInventory models.Inventory
        if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&lockedInventory, "id = ?", inventory.ID).Error; err != nil {
            return fmt.Errorf("failed to lock inventory: %w", err)
        }

        available := lockedInventory.Quantity - lockedInventory.ReservedQuantity
        if quantity > available {
            return ErrInsufficientStock
        }

        // Find and lock existing cart item
        var existing models.CartItem
        query := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("cart_id = ? AND product_id = ?", cart.ID, productID)
        if variantID != nil {
            query = query.Where("variant_id = ?", *variantID)
        } else {
            query = query.Where("variant_id IS NULL")
        }
        err := query.First(&existing).Error
        if err == nil {
            // Update existing item
            newQty := existing.Quantity + quantity
            if err := tx.Model(&models.CartItem{}).Where("id = ?", existing.ID).Update("quantity", newQty).Error; err != nil {
                return fmt.Errorf("failed to update cart item: %w", err)
            }
            // Adjust reserved quantity
            if err := tx.Model(&models.Inventory{}).Where("id = ?", inventory.ID).
                Update("reserved_quantity", gorm.Expr("reserved_quantity + ?", quantity)).Error; err != nil {
                return fmt.Errorf("failed to adjust inventory reservation: %w", err)
            }
            return nil
        } else if !errors.Is(err, gorm.ErrRecordNotFound) {
            return fmt.Errorf("failed to check existing cart item: %w", err)
        }

        // Create new cart item
        cartItem := &models.CartItem{
            CartID:     cart.ID,
            ProductID:  productID,
            VariantID:  variantID,
            Quantity:   quantity,
            MerchantID: product.MerchantID,
        }
        if err := tx.Create(cartItem).Error; err != nil {
            return fmt.Errorf("failed to create cart item: %w", err)
        }
        if err := tx.Model(&models.Inventory{}).Where("id = ?", inventory.ID).
            Update("reserved_quantity", gorm.Expr("reserved_quantity + ?", quantity)).Error; err != nil {
            return fmt.Errorf("failed to reserve inventory: %w", err)
        }
        return nil
    })
    if err != nil {
        s.logger.Error("Transaction failed", zap.Error(err))
        return nil, fmt.Errorf("%w: %v", ErrTransactionFailed, err)
    }

    var updatedCart models.Cart
    if err := db.DB.WithContext(ctx).
        Preload("CartItems.Product.Category").
        Preload("CartItems.Product.Media").
        //Preload("CartItems.Product").
        //Preload("CartItems.Variant").
        Preload("CartItems.Variant.Inventory").
        First(&updatedCart, cart.ID).Error; err != nil {
        s.logger.Error("Failed to fetch full updated cart", zap.Uint("cart_id", cart.ID), zap.Error(err))
        return nil, fmt.Errorf("failed to fetch cart: %w", err)
    }
    response := helpers.ToCartResponse(&updatedCart)
    return response, nil
}

// UpdateCartItemQuantity updates the quantity of a cart item
func (s *CartService) UpdateCartItemQuantity(ctx context.Context, cartItemID uint, quantity int) (*dto.CartResponse, error) {
	if cartItemID == 0 {
		return nil, errors.New("invalid cart item ID")
	}
	if quantity <= 0 {
		return nil, ErrInvalidQuantity
	}

	// load cart item (contains MerchantID, ProductID, and optional VariantID)
	cartItem, err := s.cartItemRepo.FindByID(ctx, cartItemID)
	if err != nil {
		return nil, repositories.ErrCartItemNotFound
	}

	// ensure we have a merchantID to scope the inventory lookup
	merchantID := cartItem.MerchantID
	if merchantID == "" {
		// fallback: fetch product to get merchant (shouldn't usually happen if cart items store merchant)
		prod, perr := s.productRepo.FindByID(ctx, cartItem.ProductID)
		if perr != nil {
			return nil, ErrInventoryNotFound
		}
		merchantID = prod.MerchantID
	}

	// Determine lookup ID: use VariantID if present, otherwise ProductID
	lookupID := cartItem.ProductID
	if cartItem.VariantID != nil {
		lookupID = *cartItem.VariantID
	}

	// Use combined repo method
	inventory, err := s.inventoryRepo.FindByProductOrVariantID(ctx, lookupID, merchantID)
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

	//return s.cartRepo.FindByID(ctx, cartItem.CartID)
    cart, err := s.cartRepo.FindByID(ctx, cartItem.CartID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error("Failed to query active cart", zap.Uint("cart_id", cartItem.CartID), zap.Error(err))
		return nil, fmt.Errorf("db error: %w", err)
	}

	// response := &dto.CartResponse{
	// 	ID:        cart.ID,
	// 	UserID:    cart.UserID,
	// 	Status:    cart.Status,
	// 	Items:     make([]dto.CartItemResponse, len(cart.CartItems)),
	// 	Total:     0,  // Will compute sum of subtotals
	// 	CreatedAt: cart.CreatedAt,
	// 	UpdatedAt: cart.UpdatedAt,
	// }
	// for i, item := range cart.CartItems {
	// 	subtotal := 0.0
	// 	// Check if Product is preloaded (avoid empty struct issues)
	// 	if item.Product.ID != "" {  // Use ID as non-zero check (struct-safe)
	// 		price := item.Product.FinalPrice.InexactFloat64()  // Convert decimal.Decimal
	// 		if item.VariantID != nil && item.Variant != nil && item.Variant.ID != "" {
	// 			price += item.Variant.FinalPrice.InexactFloat64()  // Add adjustment
	// 		}
	// 		subtotal = float64(item.Quantity) * price
	// 	}
	// 	response.Items[i] = dto.CartItemResponse{
	// 		ID:        item.ID,
	// 		ProductID: item.ProductID,
	// 		Name: item.Product.Name,
	// 		VariantID: item.VariantID,
	// 		Quantity:  item.Quantity,
	// 		Subtotal:  subtotal,  // Computed: quantity * (base + adjustment)
	// 	}
	// 	response.Total += subtotal  // Accumulate grand total
	// }
	response := helpers.ToCartResponse(cart)
	return response, nil
	
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
		prod, perr := s.productRepo.FindByID(ctx, cartItem.ProductID)
		if perr != nil {
			return nil, ErrInventoryNotFound
		}
		merchantID = prod.MerchantID
	}

	// pass ctx and merchantID as required by repo
	lookupID := cartItem.ProductID
	if cartItem.VariantID != nil {
		lookupID = *cartItem.VariantID
	}

	// Use combined repo method
	inventory, err := s.inventoryRepo.FindByProductOrVariantID(ctx, lookupID, merchantID)
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

func (s *CartService) BulkAddItems(ctx context.Context, userID uint, items dto.BulkUpdateRequest) (*dto.CartResponse, error) {
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

    // Get or create cart model (light, no full preloads)
    cartModel, err := s.getOrCreateActiveCartModel(ctx, userID) // New helper, see below
    if err != nil {
        return nil, err
    }

	err = db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        for i, item := range items.Items {
            if _, err := s.AddItemToCart(ctx, userID, item.Quantity, item.ProductID, item.VariantID); err != nil {  // Optional: if keeping single add, but better inline
                s.logger.Error("Failed to add item", zap.String("product_id", item.ProductID), zap.Stringp("variant_id", item.VariantID), zap.Error(err))
                return fmt.Errorf("failed to add item %d (product %s): %w", i+1, item.ProductID, err)
            }
            product, err := s.productRepo.FindByID(ctx, item.ProductID, "Variants.Inventory", "SimpleInventory")
            if err != nil {
                s.logger.Error("Product not found", zap.String("product_id", item.ProductID), zap.Error(err))
                return ErrProductNotFound
            }
            if product.DeletedAt.Valid {
                s.logger.Error("Product is soft-deleted", zap.String("product_id", item.ProductID))
                return ErrProductNotFound
            }

            // Determine inventory and price
            var inventory *models.Inventory
            var price decimal.Decimal = product.BasePrice
            var variantID *string = item.VariantID
            if variantID != nil && len(product.Variants) > 0 {
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
                s.logger.Error("No valid inventory", zap.String("product_id", item.ProductID), zap.Stringp("variant_id", variantID))
                return ErrInventoryNotFound
            }
            if inventory == nil {
                s.logger.Error("No valid inventory", zap.String("product_id", item.ProductID), zap.Stringp("variant_id", variantID))
                return ErrInventoryNotFound
            }

            // Transaction: Update cart item and reserve inventory
            // Lock inventory
            var lockedInventory models.Inventory
            if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&lockedInventory, "id = ?", inventory.ID).Error; err != nil {
                return fmt.Errorf("failed to lock inventory: %w", err)
            }

            available := lockedInventory.Quantity - lockedInventory.ReservedQuantity
            if item.Quantity > available {
                return ErrInsufficientStock
            }

            // Find and lock existing cart item
            var existing models.CartItem
            query := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("cart_id = ? AND product_id = ?", cartModel.ID, item.ProductID)
            if variantID != nil {
                query = query.Where("variant_id = ?", *variantID)
            } else {
                query = query.Where("variant_id IS NULL")
            }
            err = query.First(&existing).Error
            if err == nil {
                // Update existing item
                newQty := existing.Quantity + item.Quantity
                if err := tx.Model(&models.CartItem{}).Where("id = ?", existing.ID).Update("quantity", newQty).Error; err != nil {
                    return fmt.Errorf("failed to update cart item: %w", err)
                }
                // Adjust reserved quantity
                if err := tx.Model(&models.Inventory{}).Where("id = ?", inventory.ID).
                    Update("reserved_quantity", gorm.Expr("reserved_quantity + ?", item.Quantity)).Error; err != nil {
                    return fmt.Errorf("failed to adjust inventory reservation: %w", err)
                }
                // NO return nil HERE! Let the loop continue
            } else if !errors.Is(err, gorm.ErrRecordNotFound) {
                return fmt.Errorf("failed to check existing cart item: %w", err)
            } else {
                // Create new cart item
                cartItem := &models.CartItem{
                    CartID:     cartModel.ID,
                    ProductID:  item.ProductID,
                    VariantID:  variantID,
                    Quantity:   item.Quantity,
                    MerchantID: product.MerchantID,
                }
                if err := tx.Create(cartItem).Error; err != nil {
                    return fmt.Errorf("failed to create cart item: %w", err)
                }
                if err := tx.Model(&models.Inventory{}).Where("id = ?", inventory.ID).
                    Update("reserved_quantity", gorm.Expr("reserved_quantity + ?", item.Quantity)).Error; err != nil {
                    return fmt.Errorf("failed to reserve inventory: %w", err)
                }
                // NO return nil HERE either!
            }
        }
        return nil  // Success—only return AFTER the entire for loop
    })
    if err != nil {
        s.logger.Error("Transaction failed", zap.Uint("user_id", userID), zap.Error(err))
        return nil, fmt.Errorf("%w: %v", ErrTransactionFailed, err)
    }

    // Load full once at end
    fullCart, err := s.cartRepo.FindByID(ctx, cartModel.ID)
    if err != nil {
        return nil, err
    }
    return helpers.ToCartResponse(fullCart), nil
}

// New helper: Light get/create without preloads
func (s *CartService) getOrCreateActiveCartModel(ctx context.Context, userID uint) (*models.Cart, error) {
    cart, err := s.cartRepo.FindActiveCartLight(ctx, userID) // New light repo method, see below
    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, err
    }
    if errors.Is(err, gorm.ErrRecordNotFound) || cart == nil {
        newCart := &models.Cart{UserID: userID, Status: models.CartStatusActive}
        if err := s.cartRepo.Create(ctx, newCart); err != nil {
            return nil, err
        }
        cart = newCart // ID now set
    }
    return cart, nil
}