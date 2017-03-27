package gosquare

import (
	"net/http"

	"github.com/mholt/binding"
)

//WebHookEvent is the struct for a Square Web Hook Event
type WebHookEvent struct {
	MerchantID string `json:"merchant_id"`
	LocationID string `json:"location_id"`
	EventType  string `json:"event_type"`
	EntityID   string `json:"entity_id"`
}

//FieldMap is required for binding
func (obj *WebHookEvent) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&obj.MerchantID: "merchant_id",
		&obj.LocationID: "location_id",
		&obj.EventType:  "event_type",
		&obj.EntityID:   "entity_id",
	}
}
