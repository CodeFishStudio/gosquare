package gosquare

import "time"

// TokenRequest is the response for requesting a token
type TokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	AccessToken  string `json:"access_toke"`
}

// TokenResponse is the response for requesting a token
type TokenResponse struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresAt   time.Time `json:"expires_at"`
	MerchantID  string    `json:"merchant_id"`
}

//WebHookRequest is the request structs for creating a webhook
type WebHookRequest struct {
	MerchantID string   `json:"merchant_id"`
	LocationID string   `json:"location_id"`
	EventTypes []string `json:"event_type"`
	//	EntityID   string   `json:"entity_id"`
}
