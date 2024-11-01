package dto

// GenerateTokenRequest represents the token generation request
type GenerateTokenRequest struct {
	// Client name for authentication
	ClientName string `json:"client_name" binding:"required" example:"Client One Corp"`
	// Client UUID for authentication
	ClientID string `json:"client_id" binding:"required,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// GenerateTokenResponse represents the token generation response
type GenerateTokenResponse struct {
	// Indicates if the operation was successful
	Success bool `json:"success" example:"true"`
	// JWT token if successful
	Token string `json:"token,omitempty" example:"eyJhbGciOiJIUzI1NiIs..."`
	// Response message
	Message string `json:"message,omitempty" example:"Token generated successfully"`
}
