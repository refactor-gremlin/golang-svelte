package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	authapp "mysvelteapp/server_new/internal/modules/auth/app"
)

// Handlers exposes HTTP endpoints for the auth module.
// This struct bridges the HTTP layer with the application service layer,
// handling request/response transformation and HTTP-specific concerns.
//
// Data Flow Example:
// POST /auth/register
// 1. HTTP Request arrives at Handler
// 2. Request body is bound to RegisterRequest struct
// 3. Service.Register() is called with the request
// 4. Service performs business logic and database operations
// 5. AuthSuccess response is returned
// 6. Response is transformed to AuthSuccessResponse and sent as JSON
type Handlers struct {
	service *authapp.Service // Application service containing business logic
}

// NewHandlers creates a new Handlers instance with the provided auth service.
// This follows dependency injection pattern where the service dependency
// is provided from the outside, making the handlers testable.
func NewHandlers(service *authapp.Service) *Handlers {
	return &Handlers{service: service}
}

// Register godoc
// @Summary Register a new user
// @Description Creates a new user account and returns a JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Register Request"
// @Success 200 {object} AuthSuccessResponse
// @Failure 400 {object} AuthErrorResponse
// @Failure 409 {object} AuthErrorResponse
// @Router /auth/register [post]
// Register handles user registration requests.
//
// Data Flow:
// HTTP POST /auth/register
// → JSON body → RegisterRequest struct
// → Service.Register() with business logic validation
// → Database persistence via repository
// → JWT token generation
// → AuthSuccess response
// → AuthSuccessResponse JSON output
//
// Error Handling:
// - 400: Invalid JSON or validation errors
// - 409: Username/email already exists
// - 500: Internal server errors
func (h *Handlers) Register(c *gin.Context) {
	var cmd authapp.RegisterRequest
	if err := c.ShouldBindJSON(&cmd); err != nil {
		writeError(c, http.StatusBadRequest, "Invalid request payload.")
		return
	}

	result, err := h.service.Register(c.Request.Context(), cmd)
	if err != nil {
		status, message := mapAppError(err)
		writeError(c, status, message)
		return
	}

	writeSuccess(c, result)
}

// Login handles user authentication requests.
//
// Data Flow:
// HTTP POST /auth/login
// → JSON body → LoginRequest struct
// → Service.Login() with credential validation
// → Repository lookup by username
// → Password hash verification
// → JWT token generation
// → AuthSuccess response
// → AuthSuccessResponse JSON output
//
// Error Handling:
// - 400: Invalid JSON or missing fields
// - 401: Invalid credentials
// - 500: Internal server errors
func (h *Handlers) Login(c *gin.Context) {
	var cmd authapp.LoginRequest
	if err := c.ShouldBindJSON(&cmd); err != nil {
		writeError(c, http.StatusBadRequest, "Invalid request payload.")
		return
	}

	result, err := h.service.Login(c.Request.Context(), cmd)
	if err != nil {
		status, message := mapAppError(err)
		writeError(c, status, message)
		return
	}

	writeSuccess(c, result)
}

func mapAppError(err error) (int, string) {
	switch {
	case authapp.IsValidationError(err):
		return http.StatusBadRequest, err.Error()
	case authapp.IsConflictError(err):
		return http.StatusConflict, err.Error()
	case authapp.IsUnauthorizedError(err):
		return http.StatusUnauthorized, err.Error()
	default:
		return http.StatusInternalServerError, "Failed to process request."
	}
}

func writeError(c *gin.Context, status int, message string) {
	c.JSON(status, AuthErrorResponse{Message: message})
}

func writeSuccess(c *gin.Context, result *authapp.AuthSuccess) {
	c.JSON(http.StatusOK, AuthSuccessResponse{
		Token:    result.Token,
		UserID:   result.UserID,
		Username: result.Username,
	})
}
