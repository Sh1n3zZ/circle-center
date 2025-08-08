package account

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	svc "circle-center/panel/account/svc"
)

type UserHandler struct {
	userService *svc.UserService
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{
		userService: svc.NewUserService(db),
	}
}

// RegisterUser handles the user registration HTTP endpoint
// @Summary Register a new user
// @Description Register a new user account with the provided information
// @Tags account
// @Accept json
// @Produce json
// @Param request body svc.RegisterRequest true "User registration information"
// @Success 201 {object} svc.RegisterResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /account/register [post]
func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req svc.RegisterRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request data",
			Message: err.Error(),
		})
		return
	}

	// Additional validation
	if err := h.userService.ValidateRegistrationRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
		})
		return
	}

	// Register user
	response, err := h.userService.RegisterUser(c.Request.Context(), &req)
	if err != nil {
		// Handle specific errors
		switch err.Error() {
		case "username already exists":
			c.JSON(http.StatusConflict, ErrorResponse{
				Error:   "Username already exists",
				Message: "The provided username is already taken",
			})
			return
		case "email already exists":
			c.JSON(http.StatusConflict, ErrorResponse{
				Error:   "Email already exists",
				Message: "The provided email is already registered",
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Registration failed",
				Message: err.Error(),
			})
			return
		}
	}

	// Return success response
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "User registered successfully",
		"data":    response,
	})
}

// LoginUser handles the user login HTTP endpoint
// @Summary Login user
// @Description Authenticate user with email and password
// @Tags account
// @Accept json
// @Produce json
// @Param request body svc.LoginRequest true "User login information"
// @Success 200 {object} svc.LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /account/login [post]
func (h *UserHandler) LoginUser(c *gin.Context) {
	var req svc.LoginRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request data",
			Message: err.Error(),
		})
		return
	}

	// Login user
	response, err := h.userService.LoginUser(c.Request.Context(), &req)
	if err != nil {
		// Handle specific errors
		switch err.Error() {
		case "invalid email or password":
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Authentication failed",
				Message: "Invalid email or password",
			})
			return
		case "account is not active":
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Account inactive",
				Message: "Your account is not active",
			})
			return
		case "account is temporarily locked":
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Account locked",
				Message: "Your account is temporarily locked",
			})
			return
		case "account locked due to too many failed attempts":
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Account locked",
				Message: "Account locked due to too many failed attempts. Please try again later.",
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Login failed",
				Message: err.Error(),
			})
			return
		}
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successful",
		"data":    response,
	})
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// HealthCheck provides a simple health check endpoint
// @Summary Health check
// @Description Check if the account service is running
// @Tags account
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /account/health [get]
func (h *UserHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "account",
		"message": "Account service is running",
	})
}
