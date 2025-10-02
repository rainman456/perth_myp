package dispute
import (
	"context"
	"errors"
	"fmt"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrInvalidInput = errors.New("invalid input")
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
func (s *DisputeService) CreateDispute(ctx context.Context, userID uint, req dto.CreateDisputeDTO) (*dto.DisputeResponseDTO, error) {
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
func (s *DisputeService) GetDispute(ctx context.Context, disputeID string, userID uint) (*dto.DisputeResponseDTO, error) {
	dispute, err := s.disputeRepo.FindDisputeByID(ctx, disputeID)
	if err != nil {
		return nil, err
	}

	if dispute.CustomerID != userID {
		return nil, ErrUnauthorized
	}

	return mapDisputeToDTO(dispute), nil
}

func mapDisputeToDTO(d *models.Dispute) *dto.DisputeResponseDTO {
	return &dto.DisputeResponseDTO{
		ID:          d.ID,
		OrderID:     d.OrderID,
		CustomerID:  d.CustomerID,
		MerchantID:  d.MerchantID,
		Reason:      d.Reason,
		Description: d.Description,
		Status:      d.Status,
		Resolution:  d.Resolution,
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