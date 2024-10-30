package dto

// ClientQueryRequest represents the request body for client-based queries
type ClientQueryRequest struct {
    ClientName string `json:"client_name" binding:"required"`
    ClientID   string `json:"client_id" binding:"required,uuid"`
}