package models

type Product struct {
	ID    int    `json:"id"`
	SKU   string `json:"sku"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}
