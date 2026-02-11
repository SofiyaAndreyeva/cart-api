package domain

type Cart struct {
	ID    int        `json:"id"`
	Items []CartItem `json:"items"`
}
type CartItem struct {
	ID      int     `db:"id" json:"id"`
	CartID  int     `db:"cart_id" json:"cart_id"`
	Product string  `db:"product" json:"product"`
	Price   float64 `db:"price" json:"price"`
}
