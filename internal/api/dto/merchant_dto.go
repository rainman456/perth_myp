package dto


type MerchantApplyDTO struct {
    StoreName     string                 `json:"store_name" binding:"required"`
    Name          string                 `json:"name" binding:"required"`
    PersonalEmail string                 `json:"personal_email" binding:"required,email"`
    WorkEmail     string                 `json:"work_email" binding:"required,email"`
    PhoneNumber   string                 `json:"phone_number"`
    PersonalAddress map[string]any       `json:"personal_address" binding:"required"`
    WorkAddress     map[string]any       `json:"work_address" binding:"required"`
    BusinessType    string               `json:"business_type"`
    Website         string               `json:"website"`
    BusinessDescription string           `json:"business_description"`
    BusinessRegistrationNumber string    `json:"business_registration_number" binding:"required"`
    StoreLogoURL                    string `json:"store_logo_url"`
    BusinessRegistrationCertificate string `json:"business_registration_certificate"`
}



type MerchantApplyResponse struct {
    ID                   string    `json:"id"`
    StoreName            string    `json:"store_name"`
    Name                 string    `json:"name"`
    PersonalEmail        string    `json:"personal_email"`
    WorkEmail            string    `json:"work_email"`
    PhoneNumber          string    `json:"phone_number"`
    PersonalAddress      map[string]any `json:"personal_address"`
    WorkAddress          map[string]any `json:"work_address"`
    BusinessType         string         `json:"business_type"`
    Website              string         `json:"website"`
    BusinessDescription  string         `json:"business_description"`
    BusinessRegistrationNumber string   `json:"business_registration_number"`
    StoreLogoURL         string         `json:"store_logo_url"`
    BusinessRegistrationCertificate string `json:"business_registration_certificate"`
    Status               string         `json:"status"`
    CreatedAt            string         `json:"created_at"`
    UpdatedAt            string         `json:"updated_at"`
}


type MerchantLogin struct {
	WorkEmail string `json:"work_email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
}


type BankDetailsRequest struct{
	BankName    string               `json:"bank_name,omitempty"`
	Currency     string              `json:"currency,omitempty"`
	AccountName     string           `json:"account_name"`
	AccountNumber    string         `json:"account_number"`
	BankCode       string           `json:"bank_code"` 
}


