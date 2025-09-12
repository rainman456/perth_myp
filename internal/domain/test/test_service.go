package test

// import (
// 	"fmt"
// 	"api-customer-merchant/internal/db"
// 	"api-customer-merchant/internal/db/models"
// 	"api-customer-merchant/internal/db/repositories"
// 	"api-customer-merchant/internal/domain/cart"
// 	"api-customer-merchant/internal/domain/order"
// 	"api-customer-merchant/internal/domain/payment"
// 	"api-customer-merchant/internal/domain/payout"
// 	"api-customer-merchant/internal/domain/product"
// )

// func TestServices() {
// 	db.Connect()
// 	db.AutoMigrate()

// 	// Setup repositories
// 	productRepo := repositories.NewProductRepository()
// 	inventoryRepo := repositories.NewInventoryRepository()
// 	cartRepo := repositories.NewCartRepository()
// 	cartItemRepo := repositories.NewCartItemRepository()
// 	orderRepo := repositories.NewOrderRepository()
// 	orderItemRepo := repositories.NewOrderItemRepository()
// 	paymentRepo := repositories.NewPaymentRepository()
// 	payoutRepo := repositories.NewPayoutRepository()

// 	// Setup services
// 	productService := product.NewProductService(productRepo, inventoryRepo)
// 	cartService := cart.NewCartService(cartRepo, cartItemRepo, productRepo, inventoryRepo)
// 	orderService := order.NewOrderService(orderRepo, orderItemRepo, cartRepo, cartItemRepo, inventoryRepo)
// 	paymentService := payment.NewPaymentService(paymentRepo, orderRepo)
// 	payoutService := payout.NewPayoutService(payoutRepo)

// 	// Insert test data
// 	user := &models.User{Email: "test@example.com", Name: "Test User", Password: "$2a$10$examplehashedpassword", Country: "Nigeria"}
// 	if err := db.DB.Create(user).Error; err != nil {
// 		fmt.Println("Error creating user:", err)
// 	}

// 	merchant := &models.Merchant{MerchantBasicInfo: models.MerchantBasicInfo{Name: "Test Merchant", StoreName: "Test Store", PersonalEmail: "personal@example.com", WorkEmail: "work@example.com", Password: "$2a$10$examplehashedpassword"}, Status: models.MerchantStatusApproved}
// 	if err := db.DB.Create(merchant).Error; err != nil {
// 		fmt.Println("Error creating merchant:", err)
// 	}

// 	category := &models.Category{Name: "Electronics", Attributes: map[string]interface{}{"color": []string{"black"}}}
// 	if err := db.DB.Create(category).Error; err != nil {
// 		fmt.Println("Error creating category:", err)
// 	}

// 	// Test ProductService
// 	product := &models.Product{Name: "Smartphone", SKU: "SM123", Price: 599.99, CategoryID: category.ID, MerchantID: merchant.ID}
// 	if err := productService.CreateProduct(product, merchant.ID); err != nil {
// 		fmt.Println("Error creating product:", err)
// 	}
// 	p, err := productService.GetProductBySKU("SM123")
// 	if err != nil {
// 		fmt.Println("Error finding product by SKU:", err)
// 	} else {
// 		fmt.Printf("Found product by SKU: %+v\n", p)
// 	}
// 	products, err := productService.SearchProductsByName("Smart")
// 	if err != nil {
// 		fmt.Println("Error searching products:", err)
// 	} else {
// 		fmt.Printf("Found products by name: %+v\n", products)
// 	}

// 	// Add inventory
// 	inventory := &models.Inventory{ProductID: p.ID, StockQuantity: 100, LowStockThreshold: 10}
// 	if err := inventoryRepo.Create(inventory); err != nil {
// 		fmt.Println("Error creating inventory:", err)
// 	}

// 	// Test inventory stock update
// 	if err := inventoryRepo.UpdateStock(p.ID, -10); err != nil {
// 		fmt.Println("Error updating stock:", err)
// 	} else {
// 		inv, _ := inventoryRepo.FindByProductID(p.ID)
// 		fmt.Printf("Updated stock: %+v\n", inv)
// 	}

// 	// Test CartService
// 	cart, err := cartService.GetActiveCart(user.ID)
// 	if err != nil {
// 		fmt.Println("Error getting active cart:", err)
// 	} else {
// 		fmt.Printf("Active cart: %+v\n", cart)
// 	}

// 	// Add item to cart
// 	cart, err = cartService.AddItemToCart(user.ID, p.ID, 2)
// 	if err != nil {
// 		fmt.Println("Error adding item to cart:", err)
// 	} else {
// 		fmt.Printf("Cart after adding item: %+v\n", cart)
// 	}

// 	// Test OrderService
// 	order, err := orderService.CreateOrderFromCart(user.ID)
// 	if err != nil {
// 		fmt.Println("Error creating order:", err)
// 	} else {
// 		fmt.Printf("Created order: %+v\n", order)
// 	}

// 	// Verify stock after order
// 	inv, err := inventoryRepo.FindByProductID(p.ID)
// 	if err != nil {
// 		fmt.Println("Error checking stock after order:", err)
// 	} else {
// 		fmt.Printf("Stock after order: %+v\n", inv)
// 	}

// 	// Retrieve order
// 	o, err := orderService.GetOrderByID(order.ID)
// 	if err != nil {
// 		fmt.Println("Error finding order by ID:", err)
// 	} else {
// 		fmt.Printf("Found order by ID: %+v\n", o)
// 	}

// 	// Update order status
// 	o, err = orderService.UpdateOrderStatus(order.ID, string(models.OrderStatusShipped))
// 	if err != nil {
// 		fmt.Println("Error updating order status:", err)
// 	} else {
// 		fmt.Printf("Updated order status: %+v\n", o)
// 	}

// 	// Test PaymentService
// 	payment, err := paymentService.ProcessPayment(order.ID, order.TotalAmount)
// 	if err != nil {
// 		fmt.Println("Error processing payment:", err)
// 	} else {
// 		fmt.Printf("Processed payment: %+v\n", payment)
// 	}

// 	// Retrieve payment
// 	pay, err := paymentService.GetPaymentByOrderID(order.ID)
// 	if err != nil {
// 		fmt.Println("Error finding payment by order ID:", err)
// 	} else {
// 		fmt.Printf("Found payment by order ID: %+v\n", pay)
// 	}

// 	// Update payment status
// 	pay, err = paymentService.UpdatePaymentStatus(payment.ID, string(models.PaymentStatusCompleted))
// 	if err != nil {
// 		fmt.Println("Error updating payment status:", err)
// 	} else {
// 		fmt.Printf("Updated payment status: %+v\n", pay)
// 	}

// 	// Test PayoutService
// 	payout, err := payoutService.CreatePayout(merchant.ID, 500.00)
// 	if err != nil {
// 		fmt.Println("Error creating payout:", err)
// 	} else {
// 		fmt.Printf("Created payout: %+v\n", payout)
// 	}

// 	// Retrieve payout
// 	po, err := payoutService.GetPayoutByID(payout.ID)
// 	if err != nil {
// 		fmt.Println("Error finding payout by ID:", err)
// 	} else {
// 		fmt.Printf("Found payout by ID: %+v\n", po)
// 	}
// }