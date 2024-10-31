package dto

type GenerateTokenRequest struct {
	ClientName string `json:"client_name" binding:"required"`
	ClientID   string `json:"client_id" binding:"required,uuid"`
}

type GenerateTokenResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token,omitempty"`
	Message string `json:"message,omitempty"`
}
