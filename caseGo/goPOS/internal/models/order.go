package models

type CartItem struct {
	ProductID int    `json:"product_id"`
	Name      string `json:"name"`
	Price     int    `json:"price"`
	Qty       int    `json:"qty"`
}

type Order struct {
	ID       int        `json:"id"`
	Items    []CartItem `json:"items"`
	Subtotal int        `json:"subtotal"`
	Discount int        `json:"discount"`
	Total    int        `json:"total"`
	Paid     int        `json:"paid"`
	Method   string     `json:"method"`
	Change   int        `json:"change"`
	Status   string     `json:"status"`
}
