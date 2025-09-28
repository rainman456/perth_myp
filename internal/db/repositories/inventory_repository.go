package repositories

import (
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"
	"errors"

	"context"

	"gorm.io/gorm"
)

/*
type InventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository() *InventoryRepository {
	return &InventoryRepository{db: db.DB}
}

// Create adds a new inventory record
func (r *InventoryRepository) Create(inventory *models.Inventory) error {
	return r.db.Create(inventory).Error
}

// FindByProductID retrieves inventory by product ID
func (r *InventoryRepository) FindByProductID(productID string) (*models.Inventory, error) {
	var inventory models.Inventory
	err := r.db.Where("product_id = ?", productID).First(&inventory).Error
	return &inventory, err
}

// UpdateStock updates the stock quantity for a product
func (r *InventoryRepository) UpdateStock(productID string, quantityChange int) error {
	return r.db.Model(&models.Inventory{}).Where("product_id = ?", productID).
		Update("stock_quantity", gorm.Expr("stock_quantity + ?", quantityChange)).Error
}

// Update modifies an existing inventory record
func (r *InventoryRepository) Update(inventory *models.Inventory) error {
	return r.db.Save(inventory).Error
}

// Delete removes an inventory record by ID
func (r *InventoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.Inventory{}, id).Error
}

 func (r *InventoryRepository) UpdateInventoryQuantity(ctx context.Context, inventoryID uint, delta int) error { // Add this method
 	return r.db.WithContext(ctx).Model(&models.Inventory{}).Where("id = ?", inventoryID).
 		Update("quantity", gorm.Expr("quantity + ?", delta)).Error
 }
*/
type InventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository() *InventoryRepository {
	return &InventoryRepository{db: db.DB}
}

// Create adds a new vendor inventory record
func (r *InventoryRepository) Create(ctx context.Context, inv *models.Inventory) error {
	return r.db.WithContext(ctx).Create(inv).Error
}

// FindByVariantID retrieves vendor inventory by variant ID
func (r *InventoryRepository) FindByVariantID(ctx context.Context, variantID, merchantID string) (*models.Inventory, error) {
	var inv models.Inventory
	return &inv, r.db.WithContext(ctx).
		Where("variant_id = ? AND merchant_id = ?", variantID, merchantID).First(&inv).Error
}

// FindByProductID (for simple products without variants)
func (r *InventoryRepository) FindByProductID(ctx context.Context, productID string, merchantID string) (*models.Inventory, error) {
	var inv models.Inventory
	err := r.db.WithContext(ctx).
		Where("product_id = ? AND merchant_id = ?", productID, merchantID).
		First(&inv).Error
	return &inv, err
}

// UpdateStock adjusts quantity (can be negative for reservations)
func (r *InventoryRepository) UpdateStock(ctx context.Context, invID uint, delta int) error {
	return r.db.WithContext(ctx).
		Model(&models.Inventory{}).
		Where("id = ?", invID).
		Update("quantity", gorm.Expr("quantity + ?", delta)).
		Error
}

// ReserveStock increments reserved quantity
func (r *InventoryRepository) ReserveStock(ctx context.Context, invID uint, qty int) error {
	return r.db.WithContext(ctx).
		Model(&models.Inventory{}).
		Where("id = ?", invID).
		Update("reserved_quantity", gorm.Expr("reserved_quantity + ?", qty)).
		Error
}

// ReleaseStock decrements reserved quantity
func (r *InventoryRepository) ReleaseStock(ctx context.Context, invID uint, qty int) error {
	return r.db.WithContext(ctx).
		Model(&models.Inventory{}).
		Where("id = ?", invID).
		Update("reserved_quantity", gorm.Expr("reserved_quantity - ?", qty)).
		Error
}

// Delete removes a vendor inventory record by ID
func (r *InventoryRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Inventory{}, id).Error
}

// UpdateInventoryQuantity updates Quantity (can be negative)

func (r *InventoryRepository) UpdateInventoryQuantity(ctx context.Context, inventoryID string, delta int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var inv models.Inventory
		if err := tx.First(&inv, "id = ?", inventoryID).Error; err != nil {
			return err
		}
		newQ := inv.Quantity + delta
		if newQ < 0 && !inv.BackorderAllowed {
			return errors.New("insufficient stock and backorders not allowed")
		}
		inv.Quantity = newQ
		return tx.Save(&inv).Error
	})
}




// Add method for lookup by product and merchant (no VariantID)
func (r *inventoryRepository) FindByProductAndMerchant(ctx context.Context, productID, merchantID string) (*models.Inventory, error) {
	var inv models.Inventory
	err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}). // Lock for update
		Where("product_id = ? AND merchant_id = ?", productID, merchantID).
		First(&inv).Error
	return &inv, err
}

// UpdateInventory updates quantity/reserved (delta positive for unreserve)
func (r *inventoryRepository) UpdateInventory(ctx context.Context, id uint, delta int) error {
	return r.db.WithContext(ctx).
		Model(&models.Inventory{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"quantity":        gorm.Expr("quantity + ?", delta),
			"reserved_quantity": gorm.Expr("GREATEST(reserved_quantity - ?, 0)", delta),
		}).Error
}