package dto

//import "time"


type RegisterRequest struct {
    Email      string                 `json:"email" validate:"required,email"`
    Password   string                 `json:"password" validate:"required,min=6"`
    Name       string                 `json:"name" validate:"required"`
	Country  string                    `json:"country"`
    //Phone      string                 `json:"phone_number" validate:"omitempty"`
    //Address    []string      `json:"address" validate:"omitempty"`
}

type LoginResponse struct {
    Token    string `json:"token"`
    UserID   string `json:"user_id"`
    Email    string `json:"email"`
}




type LoginRequest struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}




 type ProfileResponse struct {
  	ID       uint     `json:"id"`
  	Email    string   `json:"email"`
  	Name     string   `json:"name"`
  	Country  string   `json:"country"`
  	Addresses []string `json:"addresses,omitempty"`
  }	


   type UserUpdateRequest struct {
  	Email    string   `json:"email,omitempty"`
  	Name     string   `json:"name,omitempty"`
  	Country  string   `json:"country,omitempty"`
  	Addresses []string `json:"addresses,omitempty"`
  }	


  type ResetPasswordRequest struct{
	  Email      string       `json:"email" validate:"required,email"`
	  NewPassword string `json:"new_password" binding:"required"`
  }



  type CreateAddressRequest struct {
	Address               string `json:"address" binding:"required"`
	PhoneNumber           string `json:"phone_number" binding:"omitempty"`
	AdditionalPhoneNumber string `json:"additional_phone_number" binding:"omitempty"`
	DeliveryAddress       string `json:"delivery_address" binding:"omitempty"`
	ShippingAddress       string `json:"shipping_address" binding:"omitempty"`
	State                 string `json:"state" binding:"omitempty"`
	LGA                   string `json:"lga" binding:"omitempty"`
	// Type string `json:"type" binding:"omitempty,oneof=home work billing shipping"`
}

// UpdateAddressRequest payload for updating an address
type UpdateAddressRequest struct {
	Address               *string `json:"address" binding:"omitempty"`
	PhoneNumber           *string `json:"phone_number" binding:"omitempty"`
	AdditionalPhoneNumber *string `json:"additional_phone_number" binding:"omitempty"`
	DeliveryAddress       *string `json:"delivery_address" binding:"omitempty"`
	ShippingAddress       *string `json:"shipping_address" binding:"omitempty"`
	State                 *string `json:"state" binding:"omitempty"`
	LGA                   *string `json:"lga" binding:"omitempty"`
}

// AddressResponse returned to clients
type AddressResponse struct {
	ID                    uint      `json:"id"`
	Address               string    `json:"address,omitempty"`
	PhoneNumber           string    `json:"phone_number,omitempty"`
	AdditionalPhoneNumber string    `json:"additional_phone_number,omitempty"`
	DeliveryAddress       string    `json:"delivery_address,omitempty"`
	ShippingAddress       string    `json:"shipping_address,omitempty"`
	State                 string    `json:"state,omitempty"`
	LGA                   string    `json:"lga,omitempty"`

}

// ListAddressesResponse wrapper for list responses (useful for swagger)
type ListAddressesResponse struct {
	Items []AddressResponse `json:"items"`
	Count int               `json:"count"`
}