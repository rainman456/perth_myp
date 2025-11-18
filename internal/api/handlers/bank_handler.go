package handlers

import (
	"net/http"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/bank"

	"github.com/gin-gonic/gin"
)

type BankHandler struct {
	bankService *bank.FetchBankService
}

func NewBankHandler() *BankHandler {
	return &BankHandler{
		bankService: bank.NewFetchBankService(),
	}
}

// GetBanks handles GET /banks
// @Summary List all banks
// @Description Returns list of all supported banks with codes
// @Tags Banks
// @Produce json
// @Success 200 {object} dto.BankListResponseDTO
// @Failure 500 {object} object{error=string}
// @Router /banks [get]
func (h *BankHandler) GetBanks(c *gin.Context) {
	country := c.Query("country")

	items, err := h.bankService.GetBanks(c.Request.Context(), country)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := dto.BankListResponseDTO{
		Banks: items,
		Total: len(items),
	}
	c.JSON(http.StatusOK, response)
}
