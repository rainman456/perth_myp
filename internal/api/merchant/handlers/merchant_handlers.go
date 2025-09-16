package handlers

import (
	"api-customer-merchant/internal/db/models"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"time"

	//"api-customer-merchant/internal/domain/order"
	//"api-customer-merchant/internal/domain/payout"
	"api-customer-merchant/internal/domain/product"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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

// CreateProduct godoc
// @Summary Create a new product
// @Description Allows an authenticated merchant to create a new product listing with required fields like name, SKU (unique), price, and category ID
// @Tags Merchant
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body object{name=string,sku=string,description=string,price=number,category_id=integer,stock_quantity=integer,images=array} true "Product creation details"
// @Success 201 {object} object{id=string,name=string,sku=string,description=string,price=number,category_id=integer,merchant_id=string,stock_quantity=integer,images=array,created_at=string,updated_at=string} "Product created successfully"
// @Failure 400 {object} object{error=string} "Invalid fields (e.g., missing name/SKU, invalid price, duplicate SKU)"
// @Failure 401 {object} object{error=string} "Unauthorized: Invalid or missing JWT token"
// @Failure 500 {object} object{error=string} "Failed to create product"
// @Router /merchant/create/product [post]
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

// UpdateProduct godoc
// @Summary Update an existing product
// @Description Allows an authenticated merchant to update their own product's details. SKU must remain unique
// @Tags Merchant
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Product ID (UUID)"
// @Param body body object{name=string,sku=string,description=string,price=number,category_id=integer,stock_quantity=integer,images=array} true "Updated product details"
// @Success 200 {object} object{id=string,name=string,sku=string,description=string,price=number,category_id=integer,merchant_id=string,stock_quantity=integer,images=array,created_at=string,updated_at=string} "Product updated successfully"
// @Failure 400 {object} object{error=string} "Invalid fields (e.g., duplicate SKU, invalid price)"
// @Failure 401 {object} object{error=string} "Unauthorized: Invalid or missing JWT token"
// @Failure 403 {object} object{error=string} "Forbidden: Product does not belong to this merchant"
// @Failure 500 {object} object{error=string} "Failed to update product"
// @Router /merchant/{id} [put]
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














// GetMyProducts godoc
// @Summary Retrieve all products for the authenticated merchant
// @Description Fetches a list of all products owned by the current merchant, with optional pagination and filtering
// @Tags Merchant
// @Security BearerAuth
// @Produce json
// @Param limit query integer false "Number of products to return (default: 20, max: 100)"
// @Param offset query integer false "Offset for pagination (default: 0)"
// @Param status query string false "Filter by product status (e.g., 'active', 'inactive')"
// @Param category_id query integer false "Filter by category ID"
// @Success 200 {array} models.Product "List of merchant's products"
// @Failure 401 {object} object{error=string} "Unauthorized: Invalid or missing JWT token"
// @Failure 500 {object} object{error=string} "Failed to retrieve products"
// @Router /merchant/products [get]
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
























// DeleteProduct godoc
// @Summary Delete a product owned by the merchant
// @Description Permanently deletes a product by ID, if it belongs to the authenticated merchant and has no active orders
// @Tags Merchant
// @Security BearerAuth
// @Produce json
// @Param id path string true "Product ID (UUID)"
// @Success 200 {object} object{message=string} "Product deleted successfully"
// @Failure 400 {object} object{error=string} "Invalid product ID or product not found"
// @Failure 401 {object} object{error=string} "Unauthorized: Invalid or missing JWT token"
// @Failure 403 {object} object{error=string} "Forbidden: Product does not belong to this merchant"
// @Failure 409 {object} object{error=string} "Conflict: Cannot delete product with active orders"
// @Failure 500 {object} object{error=string} "Failed to delete product"
// @Router /merchant/{id} [delete]
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
*/
























































































// BulkUploadProducts godoc
// @Summary Bulk upload products via CSV file
// @Description Allows an authenticated merchant to upload multiple products via a CSV file with header ["name","sku","price","description","category_id","stock_quantity"]
// @Tags Merchant
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param csv formData file true "CSV file for bulk upload"
// @Success 200 {object} object{message=string} "Products uploaded successfully"
// @Failure 400 {object} object{error=string} "Invalid CSV file or row validation errors"
// @Failure 401 {object} object{error=string} "Unauthorized: Invalid or missing JWT token"
// @Failure 413 {object} object{error=string} "CSV file too large"
// @Failure 500 {object} object{error=string} "Failed to process bulk upload"
// @Router /merchant/bulk-upload [post]
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