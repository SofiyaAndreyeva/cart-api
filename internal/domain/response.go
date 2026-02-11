package domain

type CartPriceResponse struct {
	CartID          int     `json:"cart_id"`
	TotalPrice      float64 `json:"total_price"`
	DiscountPercent int     `json:"discount_percent"`
	FinalPrice      float64 `json:"final_price"`
}
