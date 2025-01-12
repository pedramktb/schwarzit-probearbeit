package dtos

// @Description login request
// @Tags auth
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" format:"email" validate:"required" example:"abc@xyz.com"`
	Password string `json:"password" binding:"required" validate:"required" example:"password"`
} // @name LoginRequest

// @Description auth response
// @Tags auth
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
} // @name AuthResponse
