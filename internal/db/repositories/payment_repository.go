package repositories

import (
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

var ErrPaymentNotFound     = errors.New("payment not found")

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository() *PaymentRepository {
	return &PaymentRepository{db: db.DB}
}

// Create adds a new payment
func (r *PaymentRepository) Create(ctx context.Context,payment *models.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}


 func (r *PaymentRepository) FindByTransactionID(ctx context.Context, txID string) (*models.Payment, error) {
     var payment models.Payment
     err := r.db.WithContext(ctx).Where("transaction_id = ?", txID).First(&payment).Error
     if err != nil {
         if errors.Is(err, gorm.ErrRecordNotFound) {
             return nil, ErrPaymentNotFound
         }
         return nil, err
     }
     return &payment, nil
 }



// FindByID retrieves a payment by ID with associated Order and User
func (r *PaymentRepository) FindByID(ctx context.Context, id uint) (*models.Payment, error) {
	var payment models.Payment
	err :=  r.db.WithContext(ctx).Preload("Order.User").First(&payment, id).Error
	return &payment, err
}

// FindByOrderID retrieves a payment by order ID
func (r *PaymentRepository) FindByOrderID(ctx context.Context ,orderID uint) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.WithContext(ctx).Preload("Order.User").Where("order_id = ?", orderID).First(&payment).Error
	return &payment, err
}

// FindByUserID retrieves all payments for a user
func (r *PaymentRepository) FindByUserID(ctx context.Context ,userID uint) ([]models.Payment, error) {
	var payments []models.Payment
	err := r.db.WithContext(ctx).Preload("Order.User").Joins("JOIN orders ON orders.id = payments.order_id").Where("orders.user_id = ?", userID).Find(&payments).Error
	return payments, err
}

// Update modifies an existing payment
func (r *PaymentRepository) Update(ctx context.Context,payment *models.Payment) error {
	return r.db.WithContext(ctx).Save(payment).Error
}

// Delete removes a payment by ID
func (r *PaymentRepository) Delete(ctx context.Context,id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Payment{}, id).Error
}
