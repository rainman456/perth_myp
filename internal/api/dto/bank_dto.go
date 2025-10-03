package dto

// BankListItemDTO represents a single bank entry
type BankListItemDTO struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

// BankListResponseDTO for GET /banks
type BankListResponseDTO struct {
	Banks []BankListItemDTO `json:"banks"`
	Total int               `json:"total"`
}