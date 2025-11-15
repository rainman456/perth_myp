package repositories

import (
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"
	"context"
	"time"

	"gorm.io/gorm"
)

type DisputeRepository struct {
	db *gorm.DB
}

type ReturnRequestRepository struct {
	db *gorm.DB
}

func NewDisputeRepository() *DisputeRepository {
	return &DisputeRepository{db: db.DB}
}

func NewReturnRequestRepository() *ReturnRequestRepository {
	return &ReturnRequestRepository{db: db.DB}
}

// Create adds a new order item
func (r *DisputeRepository) Create(ctx context.Context, dispute *models.Dispute) error {
	return r.db.WithContext(ctx).Create(dispute).Error
}

func (r *DisputeRepository) FindDisputeByID(ctx context.Context, id string) (*models.Dispute, error) {
	var dispute models.Dispute
	err := r.db.WithContext(ctx).Scopes(r.activeScope()).First(&dispute, "id = ?", id).Error
	return &dispute, err
}

func (r *DisputeRepository) FindByOrderID(ctx context.Context, orderID string, customerID uint) ([]models.Dispute, error) {
	var disputes []models.Dispute
	err := r.db.WithContext(ctx).
		Scopes(r.activeScope()).
		Where("order_id = ? AND customer_id = ?", orderID, customerID).
		Preload("Order.OrderItems.Product.Category").
		Preload("Order.OrderItems.Product.Media").
		Find(&disputes).Error
	return disputes, err
}

// Update FindDisputesByCustomerID to include preloads for consistency
func (r *DisputeRepository) FindDisputesByCustomerID(ctx context.Context, customerID uint) ([]models.Dispute, error) {
	var disputes []models.Dispute
	err := r.db.WithContext(ctx).
		Scopes(r.activeScope()).
		Where("customer_id = ?", customerID).
		Preload("Order.OrderItems.Product.Category").
		Preload("Order.OrderItems.Product.Media").
		Find(&disputes).Error
	return disputes, err
}

// FindDisputesByMerchantID retrieves all disputes for a merchant
func (r *DisputeRepository) FindDisputesByMerchantID(ctx context.Context, merchantID string) ([]models.Dispute, error) {
	var disputes []models.Dispute
	err := r.db.WithContext(ctx).
		Scopes(r.activeScope()).
		Where("merchant_id = ?", merchantID).
		Preload("Order.OrderItems.Product.Category").
		Preload("Order.OrderItems.Product.Media").
		Preload("Customer").
		Find(&disputes).Error
	return disputes, err
}

// UpdateMedia updates fields
func (r *DisputeRepository) Update(ctx context.Context, dispute *models.Dispute) error {
	return r.db.WithContext(ctx).Save(dispute).Error
}

// DeleteMedia soft-deletes
func (r *DisputeRepository) DeleteDispute(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&models.Dispute{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}

// activeScope (if soft delete)
func (r *DisputeRepository) activeScope() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB { return db.Where("deleted_at IS NULL") }
}

func (r *ReturnRequestRepository) Create(ctx context.Context, returnrequests *models.ReturnRequest) error {
	return r.db.WithContext(ctx).Create(returnrequests).Error
}

// FindReturnRequestByID fetches a return request by ID
func (r *ReturnRequestRepository) FindReturnRequestByID(ctx context.Context, id string) (*models.ReturnRequest, error) {
	var returnRequest models.ReturnRequest
	err := r.db.WithContext(ctx).Scopes(r.activeScope()).First(&returnRequest, "id = ?", id).Error
	return &returnRequest, err
}

// Update updates a return request
// func (r *ReturnRequestRepository) Update(ctx context.Context, returnRequest *models.ReturnRequest) error {
//     return r.db.WithContext(ctx).Save(returnRequest).Error
// }

// DeleteReturnRequest soft-deletes a return request
func (r *ReturnRequestRepository) DeleteReturnRequest(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&models.ReturnRequest{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}

// activeScope filters out soft-deleted records
func (r *ReturnRequestRepository) activeScope() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB { return db.Where("deleted_at IS NULL") }
}

func (r *ReturnRequestRepository) FindByID(ctx context.Context, id string) (*models.ReturnRequest, error) {
	var returnRequest models.ReturnRequest
	err := r.db.WithContext(ctx).Scopes(r.activeScope()).First(&returnRequest, "id = ?", id).Error
	return &returnRequest, err
}

func (r *ReturnRequestRepository) Update(ctx context.Context, returnRequest *models.ReturnRequest) error {
	return r.db.WithContext(ctx).Save(returnRequest).Error
}

func (r *ReturnRequestRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&models.ReturnRequest{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}

// func (r *ReturnRequestRepository) activeScope() func(db *gorm.DB) *gorm.DB {
//     return func(db *gorm.DB) *gorm.DB { return db.Where("deleted_at IS NULL") }
// }

// func (r *ReturnRequestRepository) FindByCustomerID(ctx context.Context, customerID uint) ([]models.ReturnRequest, error) {
//     var returnRequests []models.ReturnRequest
//     err := r.db.WithContext(ctx).Scopes(r.activeScope()).Where("customer_id = ?", customerID).Find(&returnRequests).Error
//     return returnRequests, err
// }

func (r *ReturnRequestRepository) FindByCustomerID(ctx context.Context, customerID uint) ([]models.ReturnRequest, error) {
	var returnRequests []models.ReturnRequest
	err := r.db.WithContext(ctx).
		Scopes(r.activeScope()).
		Where("customer_id = ?", customerID).
		Preload("OrderItem.Order").
		Preload("OrderItem.Product.Category").
		Preload("OrderItem.Product.Media").
		Find(&returnRequests).Error
	return returnRequests, err
}

func (r *ReturnRequestRepository) FindByOrderID(ctx context.Context, orderID uint, customerID uint) ([]models.ReturnRequest, error) {
	var returnRequests []models.ReturnRequest
	err := r.db.WithContext(ctx).
		Scopes(r.activeScope()).
		Joins("JOIN order_items ON order_items.id = return_requests.order_item_id").
		Where("order_items.order_id = ? AND return_requests.customer_id = ?", orderID, customerID).
		Preload("OrderItem.Order").
		Preload("OrderItem.Product.Category").
		Preload("OrderItem.Product.Media").
		Find(&returnRequests).Error
	return returnRequests, err
}
