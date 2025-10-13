package repositories

import (
	"errors"

	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db"

	"gorm.io/gorm"
)

// Interface for repository (allows easy mocking)

type UserAddressRepository struct {
	db *gorm.DB
}

func NewUserAddressRepository() *UserAddressRepository {
	return &UserAddressRepository{db: db.DB}
}

func (r *UserAddressRepository) Create(addr *models.UserAddress) error {
	return r.db.Create(addr).Error
}

func (r *UserAddressRepository) GetByID(id uint) (*models.UserAddress, error) {
	var addr models.UserAddress
	if err := r.db.First(&addr, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &addr, nil
}

func (r *UserAddressRepository) ListByUser(userID uint) ([]models.UserAddress, error) {
	var addrs []models.UserAddress
	if err := r.db.Where("user_id = ?", userID).Find(&addrs).Error; err != nil {
		return nil, err
	}
	return addrs, nil
}

func (r *UserAddressRepository) Update(addr *models.UserAddress) error {
	return r.db.Save(addr).Error
}

func (r *UserAddressRepository) Delete(id uint) error {
	return r.db.Delete(&models.UserAddress{}, id).Error
}
