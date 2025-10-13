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
	response := &dto.CartResponse{
    ID:        cart.ID,
    UserID:    cart.UserID,
    Status:    cart.Status,
    Items:     make([]dto.CartItemResponse, len(cart.CartItems)),
    Total:     0,  // Will compute sum of subtotals
    CreatedAt: cart.CreatedAt,
    UpdatedAt: cart.UpdatedAt,
}
for i, item := range cart.CartItems {
    subtotal := 0.0
    // Check if Product is preloaded (avoid empty struct issues)
    if item.Product.ID != "" {  // Use ID as non-zero check (struct-safe)
        price := item.Product.FinalPrice.InexactFloat64()  // Convert decimal.Decimal
        if item.VariantID != nil && item.Variant != nil && item.Variant.ID != "" {
            price += item.Variant.FinalPrice.InexactFloat64()  // Add adjustment
        }
        subtotal = float64(item.Quantity) * price
    }
    response.Items[i] = dto.CartItemResponse{
        ID:        item.ID,
        ProductID: item.ProductID,
		Name: item.Product.Name,
        VariantID: item.VariantID,
        Quantity:  item.Quantity,
        Subtotal:  subtotal,  // Computed: quantity * (base + adjustment)
    }
    response.Total += subtotal  // Accumulate grand total
}
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
		existing, err := s.cartItemRepo.FindByProductIDAndCartID(ctx, productID, variantID, cart.ID)
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
			CartID:     cart.ID,
			ProductID:  productID,
			VariantID:  variantID, // Assume VariantID is *string in model
			Quantity:   quantity,
			MerchantID: product.MerchantID,
		}
		if err := s.cartItemRepo.Create(ctx, cartItem); err != nil {
			return fmt.Errorf("failed to create cart item: %w", err)
		}
		inventory.ReservedQuantity += quantity // Manual update since method undefined
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

	// Return updated cart
	// updatedCart, err := s.cartRepo.FindByID(ctx, cart.ID)
	// if err != nil {
	//     s.logger.Error("Failed to fetch updated cart", zap.Uint("cart_id", cart.ID), zap.Error(err))
	//     return nil, err
	// }
	// // Fix: Preload CartItems with related data
	// if err := db.DB.WithContext(ctx).
	//     Preload("CartItems.Product.Media").
	//     Preload("CartItems.Product.Variants.Inventory").
	//     Preload("CartItems.Variant").
	//     Find(updatedCart).Error; err != nil {
	//     s.logger.Error("Failed to preload cart items", zap.Error(err))
	//     return nil, err
	// }
	// return updatedCart, nil

	var updatedCart models.Cart
    if err := db.DB.WithContext(ctx).
        Preload("CartItems.Product", func(db *gorm.DB) *gorm.DB {
            return db.Select("id, base_price")
        }).
        Preload("CartItems.Variant", func(db *gorm.DB) *gorm.DB {
            return db.Select("id, price_adjustment")
        }).
        Preload("CartItems.Product.Media").
        Preload("CartItems.Product.Variants.Inventory").
        First(&updatedCart, cart.ID).Error; err != nil {
        s.logger.Error("Failed to fetch full updated cart", zap.Uint("cart_id", cart.ID), zap.Error(err))
        return nil, fmt.Errorf("failed to fetch cart: %w", err)
    }

    response := &dto.CartResponse{
        ID:        updatedCart.ID,
        UserID:    updatedCart.UserID,
        Status:    updatedCart.Status,
        Items:     make([]dto.CartItemResponse, len(updatedCart.CartItems)),
        Total:     0,
        CreatedAt: updatedCart.CreatedAt,
        UpdatedAt: updatedCart.UpdatedAt,
    }
    for i, item := range updatedCart.CartItems {
        subtotal := 0.0
        if item.Product.ID != "" {  // Struct-safe check
            price := item.Product.FinalPrice.InexactFloat64()
            if item.VariantID != nil && item.Variant != nil && item.Variant.ID != "" {
                price += item.Variant.FinalPrice.InexactFloat64()
            }
            subtotal = float64(item.Quantity) * price
        }
        response.Items[i] = dto.CartItemResponse{
            ID:        item.ID,
            ProductID: item.ProductID,
			Name: item.Product.Name,
            VariantID: item.VariantID,
            Quantity:  item.Quantity,
            Subtotal:  subtotal,
        }
        response.Total += subtotal
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
