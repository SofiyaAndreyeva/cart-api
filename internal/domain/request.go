package domain

type AddToCartRequest struct {
	Product string  `json:"product"`
	Price   float64 `json:"price"`
}
