package controllers

import (
	"encoding/json"
	"net/http"

	appauth "mysvelteapp/server/internal/application/auth"
	authmodels "mysvelteapp/server/internal/presentation/models/auth"
)

// AuthController exposes HTTP handlers for authentication routes.
type AuthController struct {
	authService *appauth.Service
}

// NewAuthController wires the controller to the auth service.
func NewAuthController(authService *appauth.Service) *AuthController {
	return &AuthController{authService: authService}
}

// Register godoc
// @Summary Register a new user
// @Description Creates a new user account and returns a JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param request body auth.RegisterRequest true "Register Request"
// @Success 200 {object} authmodels.AuthSuccessResponse
// @Failure 400 {object} authmodels.AuthErrorResponse
// @Failure 409 {object} authmodels.AuthErrorResponse
// @Router /auth/register [post]
func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var request appauth.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload.")
		return
	}

	result, err := c.authService.Register(r.Context(), request)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to process registration.")
		return
	}

	if !result.Success {
		statusCode := mapAuthErrorToStatus(result.ErrorType)
		writeJSON(w, statusCode, authmodels.AuthErrorResponse{Message: result.ErrorMessage})
		return
	}

	response := authmodels.AuthSuccessResponse{
		Token:    result.Token,
		UserID:   result.UserID,
		Username: result.Username,
	}
	writeJSON(w, http.StatusOK, response)
}

// Login godoc
// @Summary Authenticate a user
// @Description Validates credentials and returns a JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param request body auth.LoginRequest true "Login Request"
// @Success 200 {object} authmodels.AuthSuccessResponse
// @Failure 400 {object} authmodels.AuthErrorResponse
// @Failure 401 {object} authmodels.AuthErrorResponse
// @Router /auth/login [post]
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var request appauth.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload.")
		return
	}

	result, err := c.authService.Login(r.Context(), request)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to process login.")
		return
	}

	if !result.Success {
		statusCode := mapAuthErrorToStatus(result.ErrorType)
		writeJSON(w, statusCode, authmodels.AuthErrorResponse{Message: result.ErrorMessage})
		return
	}

	response := authmodels.AuthSuccessResponse{
		Token:    result.Token,
		UserID:   result.UserID,
		Username: result.Username,
	}
	writeJSON(w, http.StatusOK, response)
}

func mapAuthErrorToStatus(errorType appauth.AuthErrorType) int {
	switch errorType {
	case appauth.AuthErrorTypeValidation:
		return http.StatusBadRequest
	case appauth.AuthErrorTypeConflict:
		return http.StatusConflict
	case appauth.AuthErrorTypeUnauthorized:
		return http.StatusUnauthorized
	default:
		return http.StatusBadRequest
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, authmodels.AuthErrorResponse{Message: message})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
