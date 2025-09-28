package dto


type BankDetailsRequest struct{
	BankName    string               `json:"bank_name,omitempty"`
	Currency     string              `json:"currency,omitempty"`
	AccountName     string           `json:"account_name"`
	AccountNumber    string         `json:"account_number"`
	BankCode       string           `json:"bank_code"` 
}


