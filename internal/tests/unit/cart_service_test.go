package unit

/*
import (
	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// --- Mocks ---

type MockCartRepository struct {
	mock.Mock
}

func (m *MockCartRepository) FindActiveCart(ctx context.Context, userID uint) (*models.Cart, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Cart), args.Error(1)
}

func (m *MockCartRepository) Create(ctx context.Context, cart *models.Cart) error {
	args := m.Called(ctx, cart)
	return args.Error(0)
}

func (m *MockCartRepository) FindByID(ctx context.Context, id uint) (*models.Cart, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Cart), args.Error(1)
}

func (m *MockCartRepository) Update(ctx context.Context, cart *models.Cart) error {
	args := m.Called(ctx, cart)
	return args.Error(0)
}

type MockCartItemRepository struct {
	mock.Mock
}

func (m *MockCartItemRepository) FindByProductIDAndCartID(ctx context.Context, productID string, cartID uint) (*models.CartItem, error) {
	args := m.Called(ctx, productID, cartID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CartItem), args.Error(1)
}

func (m *MockCartItemRepository) UpdateQuantityWithReservation(ctx context.Context, itemID uint, newQuantity int, inventoryID uint) error {
	args := m.Called(ctx, itemID, newQuantity, inventoryID)
	return args.Error(0)
}

func (m *MockCartItemRepository) Create(ctx context.Context, cartItem *models.CartItem) error {
	args := m.Called(ctx, cartItem)
	return args.Error(0)
}

func (m *MockCartItemRepository) FindByID(ctx context.Context, id uint) (*models.CartItem, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CartItem), args.Error(1)
}

func (m *MockCartItemRepository) DeleteWithUnreserve(ctx context.Context, id uint, inventoryID uint) error {
	args := m.Called(ctx, id, inventoryID)
	return args.Error(0)
}

func (m *MockCartItemRepository) FindByCartID(ctx context.Context, cartID uint) ([]models.CartItem, error) {
	args := m.Called(ctx, cartID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.CartItem), args.Error(1)
}

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) FindByID(id string, preloads ...string) (*models.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

type MockInventoryRepository struct {
	mock.Mock
}

func (m *MockInventoryRepository) FindByProductID(productID string) (*models.Inventory, error) {
	args := m.Called(productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Inventory), args.Error(1)
}

func (m *MockInventoryRepository) UpdateInventoryQuantity(ctx context.Context, inventoryID uint, delta int) error {
	args := m.Called(ctx, inventoryID, delta)
	return args.Error(0)
}

// --- Test Suite ---

type CartServiceTestSuite struct {
	service       *CartService
	cartRepo      *MockCartRepository
	cartItemRepo  *MockCartItemRepository
	productRepo   *MockProductRepository
	inventoryRepo *MockInventoryRepository
	logger        *zap.Logger
}

func (suite *CartServiceTestSuite) SetupTest() {
	suite.cartRepo = new(MockCartRepository)
	suite.cartItemRepo = new(MockCartItemRepository)
	suite.productRepo = new(MockProductRepository)
	suite.inventoryRepo = new(MockInventoryRepository)
	suite.logger, _ = zap.NewDevelopment()

	// We need to cast the mocks to the interface type for the constructor
	var cartRepoInterface repositories.CartRepository = suite.cartRepo
	var cartItemRepoInterface repositories.CartItemRepository = suite.cartItemRepo
	var productRepoInterface repositories.ProductRepository = suite.productRepo
	var inventoryRepoInterface repositories.InventoryRepository = suite.inventoryRepo

	suite.service = NewCartService(
		&cartRepoInterface,
		&cartItemRepoInterface,
		&productRepoInterface,
		&inventoryRepoInterface,
		suite.logger,
	)
}

func TestCartService(t *testing.T) {
	// This is a placeholder function that testify's suite runner will use
}

func TestGetActiveCart_Existing(t *testing.T) {
	// Setup
	suite := new(CartServiceTestSuite)
	suite.SetupTest()
	ctx := context.Background()
	userID := uint(1)
	expectedCart := &models.Cart{Model: gorm.Model{ID: 1}, UserID: userID}

	suite.cartRepo.On("FindActiveCart", ctx, userID).Return(expectedCart, nil)

	// Execute
	cart, err := suite.service.GetActiveCart(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedCart, cart)
	suite.cartRepo.AssertExpectations(t)
}

func TestGetActiveCart_New(t *testing.T) {
	// Setup
	suite := new(CartServiceTestSuite)
	suite.SetupTest()
	ctx := context.Background()
	userID := uint(1)
	newCart := &models.Cart{UserID: userID, Status: models.CartStatusActive}
	createdCart := &models.Cart{Model: gorm.Model{ID: 2}, UserID: userID, Status: models.CartStatusActive}

	suite.cartRepo.On("FindActiveCart", ctx, userID).Return(nil, gorm.ErrRecordNotFound)
	suite.cartRepo.On("Create", ctx, newCart).Return(nil)
	suite.cartRepo.On("FindByID", ctx, newCart.ID).Return(createdCart, nil)

	// Execute
	cart, err := suite.service.GetActiveCart(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, createdCart, cart)
	suite.cartRepo.AssertExpectations(t)
}

func TestGetActiveCart_InvalidUser(t *testing.T) {
	// Setup
	suite := new(CartServiceTestSuite)
	suite.SetupTest()
	ctx := context.Background()

	// Execute
	_, err := suite.service.GetActiveCart(ctx, 0)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidUserID, err)
}

func TestAddItemToCart_NewItem(t *testing.T) {
	// Setup
	suite := new(CartServiceTestSuite)
	suite.SetupTest()
	ctx := context.Background()
	userID := uint(1)
	productID := "prod-123"
	quantity := uint(2)

	cart := &models.Cart{Model: gorm.Model{ID: 1}, UserID: userID}
	product := &models.Product{ID: productID, MerchantID: "merch-1"}
	inventory := &models.Inventory{Model: gorm.Model{ID: 10}, ProductID: productID, StockQuantity: 10}

	suite.cartRepo.On("FindActiveCart", ctx, userID).Return(cart, nil)
	suite.productRepo.On("FindByID", productID).Return(product, nil)
	suite.inventoryRepo.On("FindByProductID", productID).Return(inventory, nil)
	suite.cartItemRepo.On("FindByProductIDAndCartID", ctx, productID, cart.ID).Return(nil, gorm.ErrRecordNotFound)

	// Mocking the transaction is complex. We'll assume the inner calls are what matters.
	// In a real scenario with db.DB, this would be an integration test.
	// For unit tests, we assume the transaction block works if its components are mocked correctly.
	suite.cartItemRepo.On("Create", ctx, mock.AnythingOfType("*models.CartItem")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*models.CartItem)
		assert.Equal(t, cart.ID, arg.CartID)
		assert.Equal(t, productID, arg.ProductID)
		assert.Equal(t, int(quantity), arg.Quantity)
	})
	suite.inventoryRepo.On("UpdateInventoryQuantity", ctx, inventory.ID, -int(quantity)).Return(nil)
	suite.cartRepo.On("FindByID", ctx, cart.ID).Return(cart, nil)

	// Execute
	_, err := suite.service.AddItemToCart(ctx, userID, quantity, productID)

	// Assert
	assert.NoError(t, err)
	suite.cartRepo.AssertExpectations(t)
	suite.productRepo.AssertExpectations(t)
	suite.inventoryRepo.AssertExpectations(t)
	suite.cartItemRepo.AssertExpectations(t)
}

func TestAddItemToCart_ExistingItem(t *testing.T) {
	// Setup
	suite := new(CartServiceTestSuite)
	suite.SetupTest()
	ctx := context.Background()
	userID := uint(1)
	productID := "prod-123"
	quantity := uint(2)

	cart := &models.Cart{Model: gorm.Model{ID: 1}, UserID: userID}
	product := &models.Product{ID: productID, MerchantID: "merch-1"}
	inventory := &models.Inventory{Model: gorm.Model{ID: 10}, ProductID: productID, StockQuantity: 10}
	existingItem := &models.CartItem{Model: gorm.Model{ID: 100}, ProductID: productID, CartID: cart.ID, Quantity: 3}

	suite.cartRepo.On("FindActiveCart", ctx, userID).Return(cart, nil)
	suite.productRepo.On("FindByID", productID).Return(product, nil)
	suite.inventoryRepo.On("FindByProductID", productID).Return(inventory, nil)
	suite.cartItemRepo.On("FindByProductIDAndCartID", ctx, productID, cart.ID).Return(existingItem, nil)

	newQty := int(quantity) + existingItem.Quantity
	suite.cartItemRepo.On("UpdateQuantityWithReservation", ctx, existingItem.ID, newQty, inventory.ID).Return(nil)
	suite.cartRepo.On("FindByID", ctx, cart.ID).Return(cart, nil)

	// Execute
	_, err := suite.service.AddItemToCart(ctx, userID, quantity, productID)

	// Assert
	assert.NoError(t, err)
	suite.cartRepo.AssertExpectations(t)
	suite.productRepo.AssertExpectations(t)
	suite.inventoryRepo.AssertExpectations(t)
	suite.cartItemRepo.AssertExpectations(t)
}

func TestAddItemToCart_InsufficientStock(t *testing.T) {
	// Setup
	suite := new(CartServiceTestSuite)
	suite.SetupTest()
	ctx := context.Background()
	userID := uint(1)
	productID := "prod-123"
	quantity := uint(5)

	cart := &models.Cart{Model: gorm.Model{ID: 1}, UserID: userID}
	product := &models.Product{ID: productID, MerchantID: "merch-1"}
	inventory := &models.Inventory{Model: gorm.Model{ID: 10}, ProductID: productID, StockQuantity: 3}

	suite.cartRepo.On("FindActiveCart", ctx, userID).Return(cart, nil)
	suite.productRepo.On("FindByID", productID).Return(product, nil)
	suite.inventoryRepo.On("FindByProductID", productID).Return(inventory, nil)

	// Execute
	_, err := suite.service.AddItemToCart(ctx, userID, quantity, productID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrInsufficientStock, err)
	suite.inventoryRepo.AssertExpectations(t)
}

func TestUpdateCartItemQuantity_Success(t *testing.T) {
	// Setup
	suite := new(CartServiceTestSuite)
	suite.SetupTest()
	ctx := context.Background()
	cartItemID := uint(100)
	newQuantity := 5

	cartItem := &models.CartItem{Model: gorm.Model{ID: cartItemID}, ProductID: "prod-123", CartID: 1, Quantity: 2}
	inventory := &models.Inventory{Model: gorm.Model{ID: 10}, ProductID: "prod-123", StockQuantity: 10}
	cart := &models.Cart{Model: gorm.Model{ID: 1}}

	suite.cartItemRepo.On("FindByID", ctx, cartItemID).Return(cartItem, nil)
	suite.inventoryRepo.On("FindByProductID", cartItem.ProductID).Return(inventory, nil)
	suite.cartItemRepo.On("UpdateQuantityWithReservation", ctx, cartItemID, newQuantity, inventory.ID).Return(nil)
	suite.cartRepo.On("FindByID", ctx, cartItem.CartID).Return(cart, nil)

	// Execute
	_, err := suite.service.UpdateCartItemQuantity(ctx, cartItemID, newQuantity)

	// Assert
	assert.NoError(t, err)
	suite.cartItemRepo.AssertExpectations(t)
	suite.inventoryRepo.AssertExpectations(t)
	suite.cartRepo.AssertExpectations(t)
}

func TestUpdateCartItemQuantity_InvalidQuantity(t *testing.T) {
	// Setup
	suite := new(CartServiceTestSuite)
	suite.SetupTest()
	ctx := context.Background()

	// Execute
	_, err := suite.service.UpdateCartItemQuantity(ctx, 1, 0)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidQuantity, err)
}

func TestUpdateCartItemQuantity_ItemNotFound(t *testing.T) {
	// Setup
	suite := new(CartServiceTestSuite)
	suite.SetupTest()
	ctx := context.Background()
	cartItemID := uint(999)

	suite.cartItemRepo.On("FindByID", ctx, cartItemID).Return(nil, repositories.ErrCartItemNotFound)

	// Execute
	_, err := suite.service.UpdateCartItemQuantity(ctx, cartItemID, 5)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, repositories.ErrCartItemNotFound, err)
	suite.cartItemRepo.AssertExpectations(t)
}

func TestRemoveCartItem_Success(t *testing.T) {
	// Setup
	suite := new(CartServiceTestSuite)
	suite.SetupTest()
	ctx := context.Background()
	cartItemID := uint(100)

	cartItem := &models.CartItem{Model: gorm.Model{ID: cartItemID}, ProductID: "prod-123", CartID: 1, Quantity: 2}
	inventory := &models.Inventory{Model: gorm.Model{ID: 10}, ProductID: "prod-123", StockQuantity: 10}
	cart := &models.Cart{Model: gorm.Model{ID: 1}}

	suite.cartItemRepo.On("FindByID", ctx, cartItemID).Return(cartItem, nil)
	suite.inventoryRepo.On("FindByProductID", cartItem.ProductID).Return(inventory, nil)
	suite.cartItemRepo.On("DeleteWithUnreserve", ctx, cartItemID, inventory.ID).Return(nil)
	suite.cartRepo.On("FindByID", ctx, cartItem.CartID).Return(cart, nil)

	// Execute
	_, err := suite.service.RemoveCartItem(ctx, cartItemID)

	// Assert
	assert.NoError(t, err)
	suite.cartItemRepo.AssertExpectations(t)
	suite.inventoryRepo.AssertExpectations(t)
	suite.cartRepo.AssertExpectations(t)
}

func TestRemoveCartItem_ItemNotFound(t *testing.T) {
	// Setup
	suite := new(CartServiceTestSuite)
	suite.SetupTest()
	ctx := context.Background()
	cartItemID := uint(999)

	suite.cartItemRepo.On("FindByID", ctx, cartItemID).Return(nil, repositories.ErrCartItemNotFound)

	// Execute
	_, err := suite.service.RemoveCartItem(ctx, cartItemID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, repositories.ErrCartItemNotFound, err)
	suite.cartItemRepo.AssertExpectations(t)
}

func TestGetCartItemByID_Success(t *testing.T) {
	// Setup
	suite := new(CartServiceTestSuite)
	suite.SetupTest()
	ctx := context.Background()
	cartItemID := uint(100)
	expectedItem := &models.CartItem{Model: gorm.Model{ID: cartItemID}}

	suite.cartItemRepo.On("FindByID", ctx, cartItemID).Return(expectedItem, nil)

	// Execute
	item, err := suite.service.GetCartItemByID(ctx, cartItemID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedItem, item)
	suite.cartItemRepo.AssertExpectations(t)
}

func TestClearCart_Success(t *testing.T) {
	// Setup
	suite := new(CartServiceTestSuite)
	suite.SetupTest()
	ctx := context.Background()
	userID := uint(1)

	cart := &models.Cart{Model: gorm.Model{ID: 1}, UserID: userID, Status: models.CartStatusActive}
	items := []models.CartItem{
		{Model: gorm.Model{ID: 101}, ProductID: "prod-A", CartID: 1, Quantity: 1},
		{Model: gorm.Model{ID: 102}, ProductID: "prod-B", CartID: 1, Quantity: 2},
	}
	inventoryA := &models.Inventory{Model: gorm.Model{ID: 10}, ProductID: "prod-A"}
	inventoryB := &models.Inventory{Model: gorm.Model{ID: 11}, ProductID: "prod-B"}

	suite.cartRepo.On("FindActiveCart", ctx, userID).Return(cart, nil)
	suite.cartItemRepo.On("FindByCartID", ctx, cart.ID).Return(items, nil)

	// Mock the loop of RemoveCartItem calls
	suite.cartItemRepo.On("FindByID", ctx, items[0].ID).Return(&items[0], nil)
	suite.inventoryRepo.On("FindByProductID", items[0].ProductID).Return(inventoryA, nil)
	suite.cartItemRepo.On("DeleteWithUnreserve", ctx, items[0].ID, inventoryA.ID).Return(nil)
	suite.cartRepo.On("FindByID", ctx, cart.ID).Return(cart, nil)

	suite.cartItemRepo.On("FindByID", ctx, items[1].ID).Return(&items[1], nil)
	suite.inventoryRepo.On("FindByProductID", items[1].ProductID).Return(inventoryB, nil)
	suite.cartItemRepo.On("DeleteWithUnreserve", ctx, items[1].ID, inventoryB.ID).Return(nil)
	suite.cartRepo.On("FindByID", ctx, cart.ID).Return(cart, nil)

	suite.cartRepo.On("Update", ctx, cart).Return(nil)

	// Execute
	err := suite.service.ClearCart(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, models.CartStatusAbandoned, cart.Status)
	suite.cartRepo.AssertExpectations(t)
	suite.cartItemRepo.AssertExpectations(t)
}

func TestBulkAddItems_InvalidUserID(t *testing.T) {
	// Setup
	suite := new(CartServiceTestSuite)
	suite.SetupTest()
	ctx := context.Background()
	items := dto.BulkUpdateRequest{}

	// Execute
	_, err := suite.service.BulkAddItems(ctx, 0, items)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidUserID, err)
}

func TestBulkAddItems_ValidationError(t *testing.T) {
	// Setup
	suite := new(CartServiceTestSuite)
	suite.SetupTest()
	ctx := context.Background()
	userID := uint(1)
	// Create a request that will fail validation (e.g., Quantity=0)
	items := dto.BulkUpdateRequest{
		UserID: userID,
		Items: []struct {
			ProductID uint "json:\"product_id\" validate:\"required\""
			Quantity  int  "json:\"quantity\" validate:\"required,gt=0\""
		}{
			{ProductID: 1, Quantity: 0},
		},
	}

	// Execute
	_, err := suite.service.BulkAddItems(ctx, userID, items)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
}

func TestBulkAddItems_PartialFailure(t *testing.T) {
	// Setup
	suite := new(CartServiceTestSuite)
	suite.SetupTest()
	ctx := context.Background()
	userID := uint(1)
	items := dto.BulkUpdateRequest{
		UserID: userID,
		Items: []struct {
			ProductID uint "json:\"product_id\" validate:\"required\""
			Quantity  int  "json:\"quantity\" validate:\"required,gt=0\""
		}{
			{ProductID: 1, Quantity: 1},
			{ProductID: 2, Quantity: 1},
		},
	}

	// Mock AddItemToCart to fail on the second item
	// We can't easily mock the internal AddItemToCart call without refactoring.
	// A better approach for this specific test would be an integration test.
	// Here, we'll simulate the error by having a dependency fail.
	cart := &models.Cart{Model: gorm.Model{ID: 1}, UserID: userID}
	product1 := &models.Product{ID: "1", MerchantID: "merch-1"}
	inventory1 := &models.Inventory{Model: gorm.Model{ID: 10}, ProductID: "1", StockQuantity: 10}

	suite.cartRepo.On("FindActiveCart", ctx, userID).Return(cart, nil)
	// First item success
	suite.productRepo.On("FindByID", "1").Return(product1, nil)
	suite.inventoryRepo.On("FindByProductID", "1").Return(inventory1, nil)
	suite.cartItemRepo.On("FindByProductIDAndCartID", ctx, "1", cart.ID).Return(nil, gorm.ErrRecordNotFound)
	suite.cartItemRepo.On("Create", ctx, mock.Anything).Return(nil)
	suite.inventoryRepo.On("UpdateInventoryQuantity", ctx, inventory1.ID, -1).Return(nil)

	// Second item failure
	suite.productRepo.On("FindByID", "2").Return(nil, errors.New("product not found"))

	// Execute
	_, err := suite.service.BulkAddItems(ctx, userID, items)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to add item 2")
}
*/
