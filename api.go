package gosquare

import "time"

// TokenResponse is the response for requesting a token
type TokenResponse struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresAt   time.Time `json:"expires_at"`
	MerchantID  string    `json:"merchant_id"`
}

//WebHookRequest is the request structs for creating a webhook
type WebHookRequest struct {
	Topic   string `json:"topic"`
	Address string `json:"address"`
	Format  string `json:"format"`
}
