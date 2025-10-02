package handlers

import (
	"net/http"
	"sort"
	"sync"

	"api-customer-merchant/internal/api/dto"
	"api-customer-merchant/internal/bank"

	"github.com/gin-gonic/gin"
)

type BankHandler struct {
	bankService *bank.BankService
}

func NewBankHandler() *BankHandler {
	return &BankHandler{
		bankService: bank.GetBankService(),
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
	// Load banks (idempotent via sync.Once)
	if err := h.bankService.LoadBanks(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load banks"})
		return
	}

	// Concurrent processing for practice
	banksMap := h.bankService.GetAllBanks()
	
	// Convert map to slice concurrently
	banksChan := make(chan dto.BankListItemDTO, len(banksMap))
	var wg sync.WaitGroup

	for name, code := range banksMap {
		wg.Add(1)
		go func(n, c string) {
			defer wg.Done()
			banksChan <- dto.BankListItemDTO{Name: n, Code: c}
		}(name, code)
	}

	go func() {
		wg.Wait()
		close(banksChan)
	}()

	var banks []dto.BankListItemDTO
	for bank := range banksChan {
		banks = append(banks, bank)
	}

	// Sort alphabetically
	sort.Slice(banks, func(i, j int) bool {
		return banks[i].Name < banks[j].Name
	})

	response := dto.BankListResponseDTO{
		Banks: banks,
		Total: len(banks),
	}

	c.JSON(http.StatusOK, response)
}