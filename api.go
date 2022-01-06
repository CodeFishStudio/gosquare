package gosquare

import "time"

// TokenRequest is the data for requesting a token
type TokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	RefreshToken string `json:"refresh_token"`
	GrantType	 string `json:"grant_type"`
}

// TokenResponse is the response for requesting a token
type TokenResponse struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresAt   time.Time `json:"expires_at"`
	MerchantID  string    `json:"merchant_id"`
	RefreshToken string   `json:"refresh_token"`
}

// ErrorResponse is the error response for requesting a token
type ErrorResponse struct {
	Message string    `json:"message"`
	Type   string    `json:"type"`
}
