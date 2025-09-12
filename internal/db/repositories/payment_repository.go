package repositories

import (
	"api-customer-merchant/internal/db"
	"api-customer-merchant/internal/db/models"

	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository() *PaymentRepository {
	return &PaymentRepository{db: db.DB}
}

// Create adds a new payment
func (r *PaymentRepository) Create(payment *models.Payment) error {
	return r.db.Create(payment).Error
}

// FindByID retrieves a payment by ID with associated Order and User
func (r *PaymentRepository) FindByID(id uint) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Preload("Order.User").First(&payment, id).Error
	return &payment, err
}

// FindByOrderID retrieves a payment by order ID
func (r *PaymentRepository) FindByOrderID(orderID uint) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Preload("Order.User").Where("order_id = ?", orderID).First(&payment).Error
	return &payment, err
}

// FindByUserID retrieves all payments for a user
func (r *PaymentRepository) FindByUserID(userID uint) ([]models.Payment, error) {
	var payments []models.Payment
	err := r.db.Preload("Order.User").Joins("JOIN orders ON orders.id = payments.order_id").Where("orders.user_id = ?", userID).Find(&payments).Error
	return payments, err
}

// Update modifies an existing payment
func (r *PaymentRepository) Update(payment *models.Payment) error {
	return r.db.Save(payment).Error
}

// Delete removes a payment by ID
func (r *PaymentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Payment{}, id).Error
}