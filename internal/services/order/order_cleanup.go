// internal/services/order/order_cleanup.go (NEW FILE)
package order

import (
    "context"
    "time"
    
    "api-customer-merchant/internal/db/models"
    "go.uber.org/zap"
    "gorm.io/gorm"
)

// CleanupAbandonedOrders cancels orders that haven't been paid within 30 minutes
func (s *OrderService) CleanupAbandonedOrders(ctx context.Context) error {
    s.logger.Info("Starting abandoned order cleanup")
    
    cutoff := time.Now().Add(-30 * time.Minute)
    
    var abandonedOrders []models.Order
    err := s.db.WithContext(ctx).
        Preload("OrderItems").
        Where("status = ? AND created_at < ?", models.OrderStatusPending, cutoff).
        Find(&abandonedOrders).Error
    
    if err != nil {
        s.logger.Error("Failed to fetch abandoned orders", zap.Error(err))
        return err
    }

    if len(abandonedOrders) == 0 {
        s.logger.Info("No abandoned orders found")
        return nil
    }

    for _, order := range abandonedOrders {
        // Release inventory reservations
        err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
            // Release reserved inventory
            for _, item := range order.OrderItems {
                inventoryQuery := "merchant_id = ?"
                args := []interface{}{item.MerchantID}
                
                if item.VariantID != nil {
                    inventoryQuery += " AND variant_id = ?"
                    args = append(args, *item.VariantID)
                } else {
                    inventoryQuery += " AND product_id = ?"
                    args = append(args, item.ProductID)
                }

                // Release reserved quantity
                if err := tx.Model(&models.Inventory{}).
                    Where(inventoryQuery, args...).
                    Update("reserved_quantity", gorm.Expr("reserved_quantity - ?", item.Quantity)).
                    Error; err != nil {
                    return err
                }
            }

            // Cancel order
            order.Status = models.OrderStatusCancelled
            if err := tx.Save(&order).Error; err != nil {
                return err
            }

            return nil
        })

        if err != nil {
            s.logger.Error("Failed to cleanup abandoned order",
                zap.Uint("order_id", order.ID),
                zap.Error(err))
            continue
        }

        s.logger.Info("Abandoned order cleaned up",
            zap.Uint("order_id", order.ID),
            zap.Uint("user_id", order.UserID))
    }

    s.logger.Info("Abandoned order cleanup completed",
        zap.Int("cleaned_count", len(abandonedOrders)))
    
    return nil
}