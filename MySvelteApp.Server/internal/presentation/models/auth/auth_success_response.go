package authmodels

// AuthSuccessResponse matches the JSON contract expected by the frontend generator.
// @name AuthSuccessResponse
type AuthSuccessResponse struct {
	Token    string `json:"token"`
	UserID   uint   `json:"userId"`
	Username string `json:"username"`
}
