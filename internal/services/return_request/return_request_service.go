package return_request

import (
	"context"
	"errors"
	"sort"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/db/repositories"

	"github.com/google/uuid"
)


var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrInvalidInput = errors.New("invalid input")
	ErrNotFound = errors.New("not found")
)


type ReturnRequestService struct {
	repo *repositories.ReturnRequestRepository
}

func NewReturnRequestService(repo *repositories.ReturnRequestRepository) *ReturnRequestService {
	return &ReturnRequestService{repo: repo}
}

func (s *ReturnRequestService) CreateReturnRequest(ctx context.Context, userID uint, req dto.CreateReturnRequestDTO) (*dto.CreateReturnRequestResponseDTO, error) {
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

	return &dto.CreateReturnRequestResponseDTO{
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

func (s *ReturnRequestService) GetCustomerReturnRequests(ctx context.Context, userID uint) ([]dto.ReturnResponseDTO, error) {
	returnRequests, err := s.repo.FindByCustomerID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Group by OrderID
	groups := make(map[uint][]models.ReturnRequest)
	for _, req := range returnRequests {
		oid := req.OrderItem.OrderID
		groups[oid] = append(groups[oid], req)
	}

	dtos := make([]dto.ReturnResponseDTO, 0, len(groups))
	for _, reqs := range groups {
		dto := mapToReturnResponseDTO(reqs)
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



func (s *ReturnRequestService) GetReturnRequestsByOrderID(ctx context.Context, orderID uint, userID uint) (*dto.ReturnResponseDTO, error) {
	returnReqs, err := s.repo.FindByOrderID(ctx, orderID, userID)
	if err != nil {
		return nil, err
	}
	if len(returnReqs) == 0 {
		return nil, ErrNotFound // Or a custom error like ErrNoReturnRequestsFound
	}
	return mapToReturnResponseDTO(returnReqs), nil
}


// Updated mapping function
func mapToReturnResponseDTO(returnReqs []models.ReturnRequest) *dto.ReturnResponseDTO {
	if len(returnReqs) == 0 {
		return nil
	}

	first := returnReqs[0]
	resp := &dto.ReturnResponseDTO{
		OrderID:        first.OrderItem.OrderID,
		OrderCreatedAt: first.OrderItem.Order.CreatedAt,
		CustomerID:     first.CustomerID,
	}

	resp.Returns = make([]dto.ReturnItemDTO, len(returnReqs))
	for i, r := range returnReqs {
		imageURL := ""
		if len(r.OrderItem.Product.Media) > 0 {
			imageURL = r.OrderItem.Product.Media[0].URL // Assuming Media has a URL field
		}

		resp.Returns[i] = dto.ReturnItemDTO{
			ProductID:       r.OrderItem.ProductID,
			ProductName:     r.OrderItem.Product.Name,
			ProductImageURL: imageURL,
			CategorySlug:    r.OrderItem.Product.Category.CategorySlug, // Assuming Category has a Slug field
			Reason:          r.Reason,
			Status:          r.Status,
			CreatedAt:       r.CreatedAt,
		}
	}

	return resp
}


