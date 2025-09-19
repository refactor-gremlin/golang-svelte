package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	authapp "mysvelteapp/server_new/internal/modules/auth/app"
)

// Handlers exposes HTTP endpoints for the auth module.
type Handlers struct {
	service *authapp.Service
}

// NewHandlers wires the auth service into HTTP handlers.
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

	c.JSON(http.StatusOK, AuthSuccessResponse{
		Token:    result.Token,
		UserID:   result.UserID,
		Username: result.Username,
	})
}

// Login godoc
// @Summary Authenticate a user
// @Description Validates credentials and returns a JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login Request"
// @Success 200 {object} AuthSuccessResponse
// @Failure 400 {object} AuthErrorResponse
// @Failure 401 {object} AuthErrorResponse
// @Router /auth/login [post]
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

	c.JSON(http.StatusOK, AuthSuccessResponse{
		Token:    result.Token,
		UserID:   result.UserID,
		Username: result.Username,
	})
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
