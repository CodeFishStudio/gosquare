package gosquare

import "time"

//Payment is the struct for a Square Payment
type Payment struct {
	ID             string        `json:"id"`
	MerchantID     string        `json:"merchant_id"`
	CreatedAt      time.Time     `json:"created_at"`
	InclusiveTax   PaymentAmount `json:"inclusive_tax_money,omitempty"`
	NetTotal       PaymentAmount `json:"net_total_money,omitempty"`
	DiscountAmount PaymentAmount `json:"discount_money,omitempty"`
	RefundAmount   PaymentAmount `json:"refunded_money,omitempty"`
}

//PaymentAmount is the struct for a an amount with a Square Payment
type PaymentAmount struct {
	CurrencyCode string  `json:"currency_code"`
	Amount       float64 `json:"amount"`
}
