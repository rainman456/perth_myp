package dispute

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrInvalidInput = errors.New("invalid input")
	ErrNotFound     = errors.New("Not found")
)

type DisputeService struct {
	disputeRepo *repositories.DisputeRepository
	orderRepo   *repositories.OrderRepository
	logger      *zap.Logger
}

func NewDisputeService(
	disputeRepo *repositories.DisputeRepository,
	orderRepo *repositories.OrderRepository,
	logger *zap.Logger,
) *DisputeService {
	return &DisputeService{
		disputeRepo: disputeRepo,
		orderRepo:   orderRepo,
		logger:      logger,
	}
}

// CreateDispute creates a new dispute
func (s *DisputeService) CreateDispute(ctx context.Context, userID uint, req dto.CreateDisputeDTO) (*dto.CreateDisputeResponseDTO, error) {
	logger := s.logger.With(zap.String("operation", "CreateDispute"), zap.Uint("user_id", userID))

	// Verify order exists and belongs to user
	order, err := s.orderRepo.FindByID(ctx, parseUint(req.OrderID))
	if err != nil {
		logger.Error("Order not found", zap.Error(err))
		return nil, fmt.Errorf("order not found: %w", err)
	}

	if order.UserID != userID {
		return nil, ErrUnauthorized
	}

	// Create dispute model
	dispute := &models.Dispute{
		ID:          uuid.NewString(),
		OrderID:     req.OrderID,
		CustomerID:  userID,
		MerchantID:  order.OrderItems[0].MerchantID, // Assume first item's merchant
		Reason:      req.Reason,
		Description: req.Description,
		Status:      "open",
	}

	if err := s.disputeRepo.Create(ctx, dispute); err != nil {
		logger.Error("Failed to create dispute", zap.Error(err))
		return nil, err
	}

	return mapDisputeToDTO(dispute), nil
}

// GetDispute retrieves a dispute by ID
func (s *DisputeService) GetDispute(ctx context.Context, disputeID string, userID uint) (*dto.CreateDisputeResponseDTO, error) {
	dispute, err := s.disputeRepo.FindDisputeByID(ctx, disputeID)
	if err != nil {
		return nil, err
	}

	if dispute.CustomerID != userID {
		return nil, ErrUnauthorized
	}

	return mapDisputeToDTO(dispute), nil
}

func mapDisputeToDTO(d *models.Dispute) *dto.CreateDisputeResponseDTO {
	return &dto.CreateDisputeResponseDTO{
		ID:          d.ID,
		OrderID:     d.OrderID,
		CustomerID:  d.CustomerID,
		MerchantID:  d.MerchantID,
		Reason:      d.Reason,
		Description: d.Description,
		Status:        dto.PayoutStatus(d.Status),
		//Resolution:  d.Resolution,
		CreatedAt:   d.CreatedAt,
		ResolvedAt:  d.ResolvedAt,
	}
}

func parseUint(s string) uint {
	// Simplified parser (add robust error handling in production)
	var id uint
	fmt.Sscanf(s, "%d", &id)
	return id
}

func mapToDisputeResponseDTO(disputes []models.Dispute) *dto.DisputeResponseDTO {
	if len(disputes) == 0 {
		return nil
	}

	first := disputes[0]
	resp := &dto.DisputeResponseDTO{
		OrderID:        parseUint(first.OrderID),
		OrderCreatedAt: first.Order.CreatedAt,
		Status:        dto.PayoutStatus(first.Status),
		CustomerID:     first.CustomerID,
		MerchantID:     first.MerchantID,
	}

	resp.Disputes = make([]dto.DisputeItemDTO, len(disputes))
	for i, d := range disputes {
		imageURL := ""
		// Assuming first order item for simplicity; adjust if multiple items needed
		if len(d.Order.OrderItems) > 0 && len(d.Order.OrderItems[0].Product.Media) > 0 {
			imageURL = d.Order.OrderItems[0].Product.Media[0].URL // Assuming Media has a URL field
		}

		resolvedAt := &d.ResolvedAt
		if d.ResolvedAt.IsZero() {
			resolvedAt = nil
		}

		resp.Disputes[i] = dto.DisputeItemDTO{
			ProductID:       d.Order.OrderItems[0].ProductID,
			ProductName:     d.Order.OrderItems[0].Product.Name,
			ProductImageURL: imageURL,
			CategorySlug:    d.Order.OrderItems[0].Product.Category.CategorySlug, // Assuming Category has a Slug field
			Reason:          d.Reason,
			Description:     d.Description,
			Resolution:      d.Resolution,
			ResolvedAt:      resolvedAt,
			CreatedAt:       d.CreatedAt,
		}
	}

	return resp
}

// GetDisputesByOrderID retrieves disputes for a specific order
func (s *DisputeService) GetDisputesByOrderID(ctx context.Context, orderID string, userID uint) (*dto.DisputeResponseDTO, error) {
	disputes, err := s.disputeRepo.FindByOrderID(ctx, orderID, userID)
	if err != nil {
		return nil, err
	}
	if len(disputes) == 0 {
		return nil, ErrNotFound // Or a custom error like ErrNoDisputesFound
	}
	return mapToDisputeResponseDTO(disputes), nil
}

// GetCustomerDisputes retrieves all disputes for a customer, grouped by order
func (s *DisputeService) GetCustomerDisputes(ctx context.Context, userID uint) ([]dto.DisputeResponseDTO, error) {
	disputes, err := s.disputeRepo.FindDisputesByCustomerID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Group by OrderID
	groups := make(map[string][]models.Dispute)
	for _, d := range disputes {
		oid := d.OrderID
		groups[oid] = append(groups[oid], d)
	}

	dtos := make([]dto.DisputeResponseDTO, 0, len(groups))
	for _, disputeGroup := range groups {
		dto := mapToDisputeResponseDTO(disputeGroup)
		if dto != nil {
			dtos = append(dtos, *dto)
		}
	}

	// Sort by order created_at descending (most recent first)
	sort.Slice(dtos, func(i, j int) bool {
		return dtos[i].OrderCreatedAt.After(dtos[j].OrderCreatedAt)
	})

	return dtos, nil
}

// GetMerchantDisputes retrieves all disputes for a merchant, grouped by order
func (s *DisputeService) GetMerchantDisputes(ctx context.Context, merchantID string) ([]dto.DisputeResponseDTO, error) {
	disputes, err := s.disputeRepo.FindDisputesByMerchantID(ctx, merchantID)
	if err != nil {
		return nil, err
	}

	// Group by OrderID
	groups := make(map[string][]models.Dispute)
	for _, d := range disputes {
		oid := d.OrderID
		groups[oid] = append(groups[oid], d)
	}

	dtos := make([]dto.DisputeResponseDTO, 0, len(groups))
	for _, disputeGroup := range groups {
		dto := mapToDisputeResponseDTO(disputeGroup)
		if dto != nil {
			dtos = append(dtos, *dto)
		}
	}

	// Sort by order created_at descending (most recent first)
	sort.Slice(dtos, func(i, j int) bool {
		return dtos[i].OrderCreatedAt.After(dtos[j].OrderCreatedAt)
	})

	return dtos, nil
}

// UpdateDisputeStatus updates the status and resolution of a dispute
func (s *DisputeService) UpdateDisputeStatus(ctx context.Context, disputeID string, merchantID string, status string, resolution string) error {
	// Find the dispute
	dispute, err := s.disputeRepo.FindDisputeByID(ctx, disputeID)
	if err != nil {
		return err
	}

	// Verify the dispute belongs to the merchant
	if dispute.MerchantID != merchantID {
		return ErrUnauthorized
	}

	// Update the dispute
	dispute.Status = status
	dispute.Resolution = resolution
	if status == "resolved" {
		dispute.ResolvedAt = time.Now()
	}

	return s.disputeRepo.Update(ctx, dispute)
}
