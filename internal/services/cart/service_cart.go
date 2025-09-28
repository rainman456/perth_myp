package cart

/*
import (
	"context"
	"errors"
	"fmt"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ErrInvalidUserID       = errors.New("invalid user ID")
	ErrInvalidProductID    = errors.New("invalid product ID")
	ErrInvalidQuantity     = errors.New("invalid quantity")
	ErrCartNotFound        = errors.New("cart not found")
	ErrCartItemNotFound    = errors.New("cart item not found")
	ErrInsufficientStock   = errors.New("insufficient stock")
	ErrInvalidCartStatus   = errors.New("invalid cart status")
	ErrProductNotFound     = errors.New("product not found")
	ErrInventoryNotFound   = errors.New("inventory not found")
	ErrInvalidVariant      = errors.New("invalid variant")
)

type CartService struct {
	cartRepo      repositories.CartRepository
	cartItemRepo  repositories.CartItemRepository
	productRepo   repositories.ProductRepository
	inventoryRepo repositories.VendorInventoryRepository
	variantRepo   repositories.VariantRepository
	logger        *zap.Logger
	validator     *validator.Validate
}

func NewCartService(
	cartRepo repositories.CartRepository,
	cartItemRepo repositories.CartItemRepository,
	productRepo repositories.ProductRepository,
	inventoryRepo repositories.VendorInventoryRepository,
	variantRepo repositories.VariantRepository,
	logger *zap.Logger,
) *CartService {
	return &CartService{
		cartRepo:      cartRepo,
		cartItemRepo:  cartItemRepo,
		productRepo:   productRepo,
		inventoryRepo: inventoryRepo,
		variantRepo:   variantRepo,
		logger:        logger,
		validator:     validator.New(),
	}
}

// GetActiveCart fetches or creates an active cart for the user
func (s *CartService) GetActiveCart(ctx context.Context, userID uint) (*models.Cart, error) {
	if userID == 0 {
		return nil, ErrInvalidUserID
	}

	cart, err := s.cartRepo.FindActiveCart(ctx, userID)
	if err == nil {
		return cart, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Create new cart
	newCart := &models.Cart{
		UserID: userID,
		Status: models.CartStatusActive,
	}
	if err := s.cartRepo.Create(ctx, newCart); err != nil {
		return nil, err
	}

	return s.cartRepo.FindByID(ctx, newCart.ID)
}

// AddItemToCart adds or updates an item in the cart with inventory reservation
func (s *CartService) AddItemToCart(ctx context.Context, userID uint, quantity uint, productID string, variantID *string) (*models.Cart, error) {
	if userID == 0 {
		return nil, ErrInvalidUserID
	}
	if productID == "" {
		return nil, ErrInvalidProductID
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
	var priceSnapshot decimal.Decimal = product.BasePrice // Default to base

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
			freshInv, freshErr = s.inventoryRepo.FindByVariantIDTx(tx, *variantID, product.MerchantID)
		} else {
			// Re-fetch product inventory inside tx for freshness
			freshInv, freshErr = s.inventoryRepo.FindByProductIDTx(tx, productID, product.MerchantID)
		}

		if freshErr != nil {
			return freshErr
		}
		available = freshInv.Quantity - freshInv.ReservedQuantity
		if available < int(quantity) {
			return ErrInsufficientStock
		}

		// Existing item check (use tx)
		cartItem, err := s.cartItemRepo.FindByProductAndVariantTx(tx, productID, variantID, cart.ID)
		newQty := int(quantity)

		if err == nil {
			newQty += cartItem.Quantity
			if available < newQty {
				return ErrInsufficientStock
			}
			// Update (assumes repo handles delta reservation)
			return s.cartItemRepo.UpdateQuantityWithReservationTx(tx, cartItem.ID, newQty, freshInv.ID)
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			cartItem = &models.CartItem{
				CartID:        cart.ID,
				ProductID:     productID,
				VariantID:     variantID, // Nil-safe
				Quantity:      newQty,
				PriceSnapshot: priceSnapshot, // Set once
				MerchantID:    product.MerchantID,
			}
			if err := s.cartItemRepo.CreateTx(tx, cartItem); err != nil {
				return err
			}
			// Reserve (for create, delta = quantity)
			return s.inventoryRepo.ReserveStockTx(tx, freshInv.ID, newQty)
		}
		return err
	})
	if err != nil {
		s.logger.Error("Add to cart failed", zap.Error(err), zap.Uint("user_id", uint(userID)))
		return nil, err
	}

	return s.cartRepo.FindByID(ctx, cart.ID)
}

// BulkAddItems adds multiple items to the cart (MVP: loop over AddItemToCart)
func (s *CartService) BulkAddItems(ctx context.Context, userID uint, req dto.BulkUpdateRequest) (*models.Cart, error) {
	if userID == 0 {
		return nil, ErrInvalidUserID
	}

	if err := s.validator.Struct(req); err != nil {
		return nil, errors.New("validation failed: " + err.Error())
	}

	var finalCart *models.Cart
	for i, item := range req.Items {
		var variantID *string
		if item.VariantID != nil {
			variantID = item.VariantID
		}
		cart, err := s.AddItemToCart(ctx, userID, uint(item.Quantity), fmt.Sprintf("%d", item.ProductID), variantID)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed to add item %d: %v", i+1, err))
		}
		finalCart = cart
	}

	return finalCart, nil
}

// ... (other methods like UpdateCartItemQuantity, RemoveCartItem, etc., follow similar tx patterns)
// For brevity, implement similarly: Use Tx versions of repos, re-fetch inventory, check available.

*/
