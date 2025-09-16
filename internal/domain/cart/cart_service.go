package cart
/*
import (
	"errors"
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"

	"gorm.io/gorm"
)

type CartService struct {
	cartRepo      *repositories.CartRepository
	cartItemRepo  *repositories.CartItemRepository
	productRepo   *repositories.ProductRepository
	inventoryRepo *repositories.InventoryRepository
}

func NewCartService(cartRepo *repositories.CartRepository, cartItemRepo *repositories.CartItemRepository, productRepo *repositories.ProductRepository, inventoryRepo *repositories.InventoryRepository) *CartService {
	return &CartService{
		cartRepo:      cartRepo,
		cartItemRepo:  cartItemRepo,
		productRepo:   productRepo,
		inventoryRepo: inventoryRepo,
	}
}

// GetActiveCart retrieves or creates an active cart for a user
func (s *CartService) GetActiveCart(userID uint) (*models.Cart, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}

	// Try to find an active cart
	cart, err := s.cartRepo.FindActiveCart(userID)
	if err == nil {
		return cart, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Create a new active cart if none exists
	cart = &models.Cart{
		UserID: userID,
		Status: models.CartStatusActive,
	}
	if err := s.cartRepo.Create(cart); err != nil {
		return nil, err
	}
	return s.cartRepo.FindByID(cart.ID)
}

// AddItemToCart adds a product to the user's active cart
func (s *CartService) AddItemToCart(userID, productID uint, quantity int) (*models.Cart, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}
	if productID == 0 {
		return nil, errors.New("invalid product ID")
	}
	if quantity <= 0 {
		return nil, errors.New("quantity must be positive")
	}

	// Get active cart
	cart, err := s.GetActiveCart(userID)
	if err != nil {
		return nil, err
	}

	// Check if product exists
	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Check stock availability
	inventory, err := s.inventoryRepo.FindByProductID(productID)
	if err != nil {
		return nil, errors.New("inventory not found")
	}
	if inventory.StockQuantity < quantity {
		return nil, errors.New("insufficient stock")
	}

	// Check if product is already in cart
	cartItem, err := s.cartItemRepo.FindByProductIDAndCartID(productID, cart.ID)
	if err == nil {
		// Update quantity if item exists
		newQuantity := cartItem.Quantity + quantity
		if inventory.StockQuantity < newQuantity {
			return nil, errors.New("insufficient stock for updated quantity")
		}
		if err := s.cartItemRepo.UpdateQuantity(cartItem.ID, newQuantity); err != nil {
			return nil, err
		}
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new cart item
		cartItem = &models.CartItem{
			CartID:     cart.ID,
			ProductID:  productID,
			Quantity:   quantity,
			MerchantID: product.MerchantID,
		}
		if err := s.cartItemRepo.Create(cartItem); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	return s.cartRepo.FindByID(cart.ID)
}

// UpdateCartItemQuantity updates the quantity of a cart item
func (s *CartService) UpdateCartItemQuantity(cartItemID uint, quantity int) (*models.Cart, error) {
	if cartItemID == 0 {
		return nil, errors.New("invalid cart item ID")
	}
	if quantity <= 0 {
		return nil, errors.New("quantity must be positive")
	}

	// Get cart item
	cartItem, err := s.cartItemRepo.FindByID(cartItemID)
	if err != nil {
		return nil, errors.New("cart item not found")
	}

	// Check stock availability
	inventory, err := s.inventoryRepo.FindByProductID(cartItem.ProductID)
	if err != nil {
		return nil, errors.New("inventory not found")
	}
	if inventory.StockQuantity < quantity {
		return nil, errors.New("insufficient stock")
	}

	// Update quantity
	if err := s.cartItemRepo.UpdateQuantity(cartItemID, quantity); err != nil {
		return nil, err
	}

	return s.cartRepo.FindByID(cartItem.CartID)
}

// RemoveCartItem removes an item from the cart
func (s *CartService) RemoveCartItem(cartItemID uint) (*models.Cart, error) {
	if cartItemID == 0 {
		return nil, errors.New("invalid cart item ID")
	}

	// Get cart item to find cart ID
	cartItem, err := s.cartItemRepo.FindByID(cartItemID)
	if err != nil {
		return nil, errors.New("cart item not found")
	}

	// Delete cart item
	if err := s.cartItemRepo.Delete(cartItemID); err != nil {
		return nil, err
	}

	return s.cartRepo.FindByID(cartItem.CartID)
}

func (s *CartService) GetCartItemByID(cartItemID uint) (*models.CartItem, error) {
	if cartItemID == 0 {
		return nil, errors.New("invalid cart item ID")
	}
	return s.cartItemRepo.FindByID(cartItemID)
}
*/