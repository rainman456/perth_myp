package handlers

import (
	"log"
	"net/http"
	"os"

	//"os"
	"strings"

	//"api-customer-merchant/internal/db/models"
	//"api-customer-merchant/internal/db/repositories"
	services "api-customer-merchant/internal/domain/user"
	"api-customer-merchant/internal/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	service *services.AuthService
}

// In customer/handlers/auth_handler.go AND merchant/handlers/auth_handler.go
func NewAuthHandler(s *services.AuthService) *AuthHandler {
    return &AuthHandler{
        service: s,
    }
}
// Register godoc
// @Summary Register a new customer
// @Description Creates a new customer account with email, name, password, and optional country
// @Tags Customer
// @Accept json
// @Produce json
// @Param body body object{email=string,name=string,password=string,country=string} true "Customer registration details"
// @Success 200 {object} object{token=string} "JWT token"
// @Failure 400 {object} object{error=string} "Invalid request"
// @Failure 500 {object} object{error=string} "Server error"
// @Router /customer/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
		Country  string `json:"country"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.RegisterUser(req.Email, req.Name, req.Password, req.Country)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": token})
}


// Login godoc
// @Summary Customer login
// @Description Authenticates a customer using email and password
// @Tags Customer
// @Accept json
// @Produce json
// @Param body body object{email=string,password=string} true "Customer login credentials"
// @Success 200 {object} object{token=string} "JWT token"
// @Failure 400 {object} object{error=string} "Invalid request"
// @Failure 401 {object} object{error=string} "Unauthorized"
// @Failure 403 {object} object{error=string} "Invalid role"
// @Failure 500 {object} object{error=string} "Server error"
// @Router /customer/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.LoginUser(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// if user.Role != "customer" {
	// 	c.JSON(http.StatusForbidden, gin.H{"error": "Invalid role for this API"})
	// 	return
	// }

	token, err := h.service.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}


// GoogleAuth godoc
// @Summary Initiate Google OAuth for customer
// @Description Redirects to Google OAuth login page
// @Tags Customer
// @Produce json
// @Success 307 {object} object{} "Redirect to Google OAuth"
// @Router /customer/auth/google [get]
func (h *AuthHandler) GoogleAuth(c *gin.Context) {
	url := h.service.GetOAuthConfig("customer").AuthCodeURL("state-customer", oauth2.AccessTypeOffline)
	//c.Redirect(http.StatusTemporaryRedirect, url)
	 c.JSON(http.StatusOK, gin.H{"url": url})
}



// GoogleCallback godoc
// @Summary Handle Google OAuth callback for customer
// @Description Processes Google OAuth callback and returns JWT token
// @Tags Customer
// @Produce json
// @Param code query string true "OAuth code"
// @Success 200 {object} object{token=string} "JWT token"
// @Failure 400 {object} object{error=string} "Code not provided"
// @Failure 500 {object} object{error=string} "Server error"
// @Router /customer/auth/google/callback [get]
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
    code := c.Query("code")
    state := c.Query("state")
    if code == "" || state == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code or state"})
        return
    }
    // Verify state (in production, check against stored value)
    if state != "state-customer" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state parameter"})
        return
    }

    user, token, err := h.service.GoogleLogin(code, os.Getenv("BASE_URL"), "customer")
    if err != nil {
        log.Printf("Google login failed: %v", err)
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"token": token, "user": user})
}



// Logout godoc
// @Summary Customer logout
// @Description Invalidates the customer's JWT token
// @Tags Customer
// @Security BearerAuth
// @Produce json
// @Success 200 {object} object{message=string} "Logout successful"
// @Failure 400 {object} object{error=string} "Authorization header required"
// @Router /customer/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header required"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	utils.Add(tokenString)
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
    userID, _ := c.Get("userID")
    var req struct { Name string; Country string; Addresses []string }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := h.service.UpdateProfile(userID.(uint), req.Name, req.Country, req.Addresses); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "updated"})
}