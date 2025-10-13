package user

import (
	"context"
	"errors"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"
)

// AddressService provides address business logic
type AddressService struct {
	repo *repositories.UserAddressRepository
}

func NewAddressService(repo *repositories.UserAddressRepository) *AddressService {
	return &AddressService{repo: repo}
}

// CreateAddress creates an address for a user
func (s *AddressService) CreateAddress(ctx context.Context, userID uint, req dto.CreateAddressRequest) (*models.UserAddress, error) {
	addr := &models.UserAddress{
		UserID:                userID,
		PhoneNumber:           req.PhoneNumber,
		AdditionalPhoneNumber: req.AdditionalPhoneNumber,
		DeliveryAddress:       req.DeliveryAddress,
		ShippingAddress:       req.ShippingAddress,
		State:                 req.State,
		LGA:                   req.LGA,
	}
	if err := s.repo.Create(addr); err != nil {
		return nil, err
	}
	return addr, nil
}

// GetAddress returns a single address by id
func (s *AddressService) GetAddress(ctx context.Context, id uint) (*models.UserAddress, error) {
	return s.repo.GetByID(id)
}

// ListAddresses returns addresses for a user
func (s *AddressService) ListAddresses(ctx context.Context, userID uint) ([]models.UserAddress, error) {
	return s.repo.ListByUser(userID)
}

// UpdateAddress updates an existing address. It enforces that the address belongs to userID.
func (s *AddressService) UpdateAddress(ctx context.Context, userID uint, id uint, req dto.UpdateAddressRequest) (*models.UserAddress, error) {
	addr, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if addr == nil {
		return nil, nil
	}
	if addr.UserID != userID {
		return nil, errors.New("forbidden")
	}

	// apply updates if provided
	// if req.Address != nil {
	// 	addr.Address = *req.Address
	// }
	if req.PhoneNumber != nil {
		addr.PhoneNumber = *req.PhoneNumber
	}
	if req.AdditionalPhoneNumber != nil {
		addr.AdditionalPhoneNumber = *req.AdditionalPhoneNumber
	}
	if req.DeliveryAddress != nil {
		addr.DeliveryAddress = *req.DeliveryAddress
	}
	if req.ShippingAddress != nil {
		addr.ShippingAddress = *req.ShippingAddress
	}
	if req.State != nil {
		addr.State = *req.State
	}
	if req.LGA != nil {
		addr.LGA = *req.LGA
	}

	if err := s.repo.Update(addr); err != nil {
		return nil, err
	}
	return addr, nil
}

// DeleteAddress deletes an address, only if it belongs to userID
func (s *AddressService) DeleteAddress(ctx context.Context, userID uint, id uint) error {
	addr, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if addr == nil {
		return nil // not found, caller can interpret nil as not found
	}
	if addr.UserID != userID {
		return errors.New("forbidden")
	}
	return s.repo.Delete(id)
}
