package handlers

// import (
// 	"api-customer-merchant/internal/db/models"
// 	"encoding/csv"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"time"

// 	//"api-customer-merchant/internal/domain/order"
// 	//"api-customer-merchant/internal/domain/payout"
// 	"api-customer-merchant/internal/services/product"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// )

/*
type MerchantHandlers struct {
	productService *product.ProductService
	orderService   *order.OrderService
	payoutService  *payout.PayoutService
}

func NewMerchantHandlers(productService *product.ProductService, orderService *order.OrderService, payoutService *payout.PayoutService) *MerchantHandlers {
	return &MerchantHandlers{
		productService: productService,
		orderService:   orderService,
		payoutService:  payoutService,
	}
}

// CreateProduct handles POST /merchant/products
func (h *MerchantHandlers) CreateProduct(c *gin.Context) {
	merchantID, exists := c.Get("merchantID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.productService.CreateProduct(&product, merchantID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// UpdateProduct handles PUT /merchant/products/:id
func (h *MerchantHandlers) UpdateProduct(c *gin.Context) {
	merchantID, exists := c.Get("merchantID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product.ID = uint(id)

	if err := h.productService.UpdateProduct(&product, merchantID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct handles DELETE /merchant/products/:id
func (h *MerchantHandlers) DeleteProduct(c *gin.Context) {
	merchantID, exists := c.Get("merchantID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	if err := h.productService.DeleteProduct(uint(id), merchantID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product deleted"})
}

// GetOrders handles GET /merchant/orders
func (h *MerchantHandlers) GetOrders(c *gin.Context) {
	merchantID, exists := c.Get("merchantID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	orders, err := h.orderService.GetOrdersByMerchantID(merchantID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetPayouts handles GET /merchant/payouts
func (h *MerchantHandlers) GetPayouts(c *gin.Context) {
	merchantID, exists := c.Get("merchantID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	payouts, err := h.payoutService.GetPayoutsByMerchantID(merchantID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payouts)
}
*/

/*
type MerchantHandlers struct {
	productService *product.ProductService
	//orderService   *order.OrderService
	//payoutService  *payout.PayoutService
	//promotionService *promotions.PromotionService
}
//, orderService *order.OrderService, payoutService *payout.PayoutService,promotionService *promotions.PromotionService) *MerchantHandlers {
//orderService:   orderService,
		//payoutService:  payoutService,
		//promotionService:promotionService,
func NewMerchantHandlers(productService *product.ProductService)*MerchantHandlers{
	return &MerchantHandlers{
		productService: productService,
	}
}

func (h *MerchantHandlers) CreateProduct(c *gin.Context) {
	merchantID, exists := c.Get("merchantID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	merchantIDStr := merchantID.(string)

	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.productService.CreateProduct(&product, merchantIDStr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}


 func (h *MerchantHandlers) UpdateProduct(c *gin.Context) {
 	merchantID, exists := c.Get("merchantID")
 	if !exists {
 		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
 		return
 	}
	merchantIDStr := merchantID.(string)

 	idStr := c.Param("id")
 	//id, err := strconv.ParseUint(idStr, 10, 32)
 	// if err != nil {
 	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
 	// 	return
 	// }

	var product models.Product
 	if err := c.ShouldBindJSON(&product); err != nil {
 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
 		return
 	}
 	product.ID = idStr

 	if err := h.productService.UpdateProduct(&product, merchantIDStr); err != nil {
 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
 		return
 	}

 	c.JSON(http.StatusOK, product)
 }







// func (h *MerchantHandler) UpdateProduct(c *gin.Context) {
// 	merchantID, exists := c.Get("merchantID")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
// 		return
// 	}
// 	merchantIDStr := merchantID.(string)

// 	id := c.Param("id")

// 	var product models.Product
// 	if err := c.ShouldBindJSON(&product); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	product.ID = id
// 	product.UpdatedAt = time.Now()

// 	if err := h.productService.UpdateProduct(&product, merchantIDStr); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, product)
// }



func (h *MerchantHandlers) GetMyProducts(c *gin.Context) {
	merchantID, exists := c.Get("merchantID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	merchantIDStr := merchantID.(string)

	//products, err := h.productService.GetProductsByMerchantID(merchantIDStr)
	products, err := h.productService.GetProductsByMerchantID(merchantIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}


	c.JSON(http.StatusOK, products)
}

























func (h *MerchantHandlers) DeleteProduct(c *gin.Context) {
	merchantID, exists := c.Get("merchantID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	merchantIDStr := merchantID.(string)


	idStr := c.Param("id")
	///id, err := strconv.ParseUint(idStr, 10, 32)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
	// 	return
	// }
	id:=idStr

	if err := h.productService.DeleteProduct(id,merchantIDStr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product deleted"})
}
/*
// GetOrders handles GET /merchant/orders
func (h *MerchantHandlers) GetOrders(c *gin.Context) {
	merchantID, exists := c.Get("merchantID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	orders, err := h.orderService.GetOrdersByMerchantID(merchantID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetPayouts handles GET /merchant/payouts
func (h *MerchantHandlers) GetPayouts(c *gin.Context) {
	merchantID, exists := c.Get("merchantID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	payouts, err := h.payoutService.GetPayoutsByMerchantID(merchantID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payouts)


}


func (h *MerchantHandlers) CreatePromotion(c *gin.Context) {
    merchantID, _ := c.Get("merchantID")
    var promo models.Promotion
    if err := c.ShouldBindJSON(&promo); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    promo.MerchantID = merchantID.(uint)
    if err := h.promotionService.CreatePromotion(&promo); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, promo)
}

func (h *MerchantHandlers) UpdateOrderItemStatus(c *gin.Context) {
    merchantID, _ := c.Get("merchantID")
    itemIDStr := c.Param("itemID")
    itemID, _ := strconv.ParseUint(itemIDStr, 10, 32)
    var req struct { Status string `json:"status"` }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := h.orderService.UpdateOrderItemStatus(uint(itemID), models.FulfillmentStatus(req.Status), merchantID.(uint)); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "updated"})
}







 func (h *MerchantHandlers) BulkUploadProducts(c *gin.Context) {
     merchantID, exists := c.Get("merchantID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	merchantIDStr := merchantID.(string)
     file, err := c.FormFile("csv")
     if err != nil {
         c.JSON(http.StatusBadRequest, gin.H{"error": "csv file required"})
         return
     }
     f, err := file.Open()
     if err != nil {
         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
         return
     }
     defer f.Close()

     reader := csv.NewReader(f)
     // Skip header
     _, err = reader.Read()
     if err != nil {
         c.JSON(http.StatusBadRequest, gin.H{"error": "invalid csv header"})
         return
     }

     for {
         record, err := reader.Read()
        if err == io.EOF {
             break
         }
         if err != nil {
             c.JSON(http.StatusBadRequest, gin.H{"error": "invalid csv format"})
             return
         }

//         // Parse price with error handling
         price, err := strconv.ParseFloat(record[2], 64)
         if err != nil {
             c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid price format in CSV: %s", record[2])})
             return
         }

//         // Parse record: name, sku, price, etc.
         product := &models.Product{
             Name:       record[0],
             SKU:        record[1],
             Price:      price,
             // Add other fields as needed (e.g., CategoryID if required in CSV)
			 ID:        uuid.New().String(),
			CreatedAt: time.Now(),
             MerchantID:merchantIDStr,
         }

         if err := h.productService.CreateProduct(product, merchantIDStr); err != nil {
             c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to create product: %s", err.Error())})
             return
         }
     }
     c.JSON(http.StatusOK, gin.H{"message": "products uploaded"})
 }
*/

/*
func (h *MerchantHandler) BulkUploadProducts(c *gin.Context) {
	merchantID, exists := c.Get("merchantID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	merchantIDStr := merchantID.(string)

	file, err := c.FormFile("csv")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "csv file required"})
		return
	}
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer f.Close()

	reader := csv.NewReader(f)
	_, err = reader.Read() // Skip header
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid csv header"})
		return
	}

	successCount := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid csv format"})
			return
		}

		price, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid price: %s", record[2])})
			return
		}

		product := &models.Product{
			Name:      record[0],
			SKU:       record[1],
			Price:     price,
			MerchantID: merchantIDStr,
			ID:        uuid.New().String(),
			CreatedAt: time.Now(),
		}

		if err := h.productService.CreateProduct(&product, merchantIDStr); err != nil {
			log.Printf("Failed to upload product %s: %v", record[0], err)
			continue // Continue on error (or abort if strict)
		}
		successCount++
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("uploaded %d products", successCount)})
}
*/



// func AddBankDetails(c *gin.Context) {
//     merchantID := c.GetString("merchantID")  // From auth middleware
//     var req dto.BankDetailsRequest  // Define DTO with fields
//     if err := c.ShouldBindJSON(&req); err != nil { /* handle */ }

//     err := merchantService.AddBankDetails(merchantID, MerchantBankDetails{
//         // Map fields
//     })
//     if err != nil { /* handle */ }
//     c.JSON(200, gin.H{"message": "Bank details added"})
// }