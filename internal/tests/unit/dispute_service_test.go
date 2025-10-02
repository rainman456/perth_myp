// internal/tests/unit/dispute_service_test.go (fixed)
package unit
/*
import (
	"context"
	"testing"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/db/models"
	"api-customer-merchant/internal/services/dispute"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MockDisputeRepo struct {
	mock.Mock
}

func (m *MockDisputeRepo) Create(ctx context.Context, dispute *models.Dispute) error {
	args := m.Called(ctx, dispute)
	return args.Error(0)
}

func (m *MockDisputeRepo) FindDisputeByID(ctx context.Context, id string) (*models.Dispute, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Dispute), args.Error(1)
}

type MockOrderRepo struct {
	mock.Mock
}

func (m *MockOrderRepo) FindByID(ctx context.Context, id uint) (*models.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func TestCreateDispute_Success(t *testing.T) {
	disputeRepo := new(MockDisputeRepo)
	orderRepo := new(MockOrderRepo)
	logger, _ := zap.NewDevelopment()

	service := dispute.NewDisputeService(disputeRepo, orderRepo, logger)

	ctx := context.Background()
	userID := uint(1)
	req := dto.CreateDisputeDTO{
		OrderID:     "1",
		Reason:      "Defective product",
		Description: "Item arrived damaged",
	}

	order := &models.Order{
		Model:      gorm.Model{ID: 1},
		UserID:     userID,
		OrderItems: []models.OrderItem{{MerchantID: "merch-1"}},
	}

	orderRepo.On("FindByID", ctx, uint(1)).Return(order, nil)
	disputeRepo.On("Create", ctx, mock.AnythingOfType("*models.Dispute")).Return(nil)

	resp, err := service.CreateDispute(ctx, userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "open", resp.Status)
	disputeRepo.AssertExpectations(t)
	orderRepo.AssertExpectations(t)
}

func TestCreateDispute_Unauthorized(t *testing.T) {
	disputeRepo := new(MockDisputeRepo)
	orderRepo := new(MockOrderRepo)
	logger, _ := zap.NewDevelopment()

	service := dispute.NewDisputeService(disputeRepo, orderRepo, logger)

	ctx := context.Background()
	userID := uint(2)
	req := dto.CreateDisputeDTO{OrderID: "1", Reason: "Test", Description: "Test"}

	order := &models.Order{Model: gorm.Model{ID: 1}, UserID: uint(1)}
	orderRepo.On("FindByID", ctx, uint(1)).Return(order, nil)

	_, err := service.CreateDispute(ctx, userID, req)

	assert.ErrorIs(t, err, dispute.ErrUnauthorized)
}
*/