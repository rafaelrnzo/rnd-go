package goPOS

type Product struct {
	Id      int64
	Product string
	Type    string
	Price   float64
}

type Order struct {
	ProductId int64
	Quantity  int64
}

type User struct {
	Name    string
	Email   string
	History string
}
