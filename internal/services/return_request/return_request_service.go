package return_request

import (
	"context"
	"errors"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"

	"github.com/google/uuid"
)


var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrInvalidInput = errors.New("invalid input")
)


type ReturnRequestService struct {
	repo *repositories.ReturnRequestRepository
}

func NewReturnRequestService(repo *repositories.ReturnRequestRepository) *ReturnRequestService {
	return &ReturnRequestService{repo: repo}
}

func (s *ReturnRequestService) CreateReturnRequest(ctx context.Context, userID uint, req dto.CreateReturnRequestDTO) (*dto.ReturnRequestResponseDTO, error) {
	returnReq := &models.ReturnRequest{
		ID:          uuid.NewString(),
		OrderItemID: req.OrderItemID,
		CustomerID:  userID,
		Reason:      req.Reason,
		Status:      "Pending",
	}

	if err := s.repo.Create(ctx, returnReq); err != nil {
		return nil, err
	}

	return &dto.ReturnRequestResponseDTO{
		ID:          returnReq.ID,
		OrderItemID: returnReq.OrderItemID,
		CustomerID:  returnReq.CustomerID,
		Reason:      returnReq.Reason,
		Status:      returnReq.Status,
		CreatedAt:   returnReq.CreatedAt,
		UpdatedAt:   returnReq.UpdatedAt,
	}, nil
}

func (s *ReturnRequestService) GetReturnRequest(ctx context.Context, id string, userID uint) (*dto.ReturnRequestResponseDTO, error) {
    returnReq, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }

    if returnReq.CustomerID != userID {
        return nil, ErrUnauthorized
    }

    return mapReturnRequestToDTO(returnReq), nil
}

func (s *ReturnRequestService) GetCustomerReturnRequests(ctx context.Context, userID uint) ([]dto.ReturnRequestResponseDTO, error) {
    returnRequests, err := s.repo.FindByCustomerID(ctx, userID)
    if err != nil {
        return nil, err
    }

    dtos := make([]dto.ReturnRequestResponseDTO, len(returnRequests))
    for i, req := range returnRequests {
        dtos[i] = *mapReturnRequestToDTO(&req)
    }
    return dtos, nil
}

func mapReturnRequestToDTO(r *models.ReturnRequest) *dto.ReturnRequestResponseDTO {
    return &dto.ReturnRequestResponseDTO{
        ID:          r.ID,
        OrderItemID: r.OrderItemID,
        CustomerID:  r.CustomerID,
        Reason:      r.Reason,
        Status:      r.Status,
        CreatedAt:   r.CreatedAt,
        UpdatedAt:   r.UpdatedAt,
    }
}