package authmodels

// AuthErrorResponse wraps error messages in a serialisable structure.
// @name AuthErrorResponse
type AuthErrorResponse struct {
	Message string `json:"message"`
}
