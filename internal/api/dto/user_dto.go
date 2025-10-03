package dto


type RegisterRequest struct {
    Email      string                 `json:"email" validate:"required,email"`
    Password   string                 `json:"password" validate:"required,min=6"`
    Name       string                 `json:"name" validate:"required"`
	Country  string                    `json:"country"`
    Phone      string                 `json:"phone_number" validate:"omitempty"`
    Address    string      `json:"address" validate:"omitempty"`
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
