package services

import (
	"errors"
	"fmt"
	"goPOS/internal/models"
	"goPOS/internal/storage"
)

func NewOrder() *models.Order {
	o := models.Order{
		ID:     storage.NextOrderID,
		Status: "DRAFT",
	}
	storage.NextOrderID++
	storage.Orders = append(storage.Orders, o)
	return &storage.Orders[len(storage.Orders)-1]
}

func GetOrderByID(id int) (*models.Order, error) {
	for i := range storage.Orders {
		if storage.Orders[i].ID == id {
			return &storage.Orders[i], nil
		}
	}
	return nil, fmt.Errorf("Order with ID %d not found", id)
}

func AddItem(order *models.Order, product *models.Product, qty int) error {
	if order.Status != "DRAFT" {
		return fmt.Errorf("Order with ID %d already closed", order.ID)
	}

	if qty <= 0 {
		return fmt.Errorf("Quantity must be greater than zero")
	}

	for i := range order.Items {
		if order.Items[i].ProductID == product.ID {
			order.Items[i].Qty += qty
			return fmt.Errorf("Order with ID %d already exists", order.ID)
			recalc(order)
			return nil
		}
	}
	order.Items = append(order.Items, models.CartItem{
		ProductID: product.ID,
		Name:      product.Name,
		Price:     product.Price,
		Qty:       qty,
	})
	recalc(order)
	return nil
}

func RemoveItem(order *models.Order, productID int) error {
	if order.Status != "DRAFT" {
		return fmt.Errorf("Order with ID %d already closed", order.ID)
	}
	found := false
	newItems := make([]models.CartItem, 0, len(order.Items))
	for _, it := range order.Items {
		if it.ProductID == productID {
			found = true
			continue
		}
		newItems = append(newItems, it)
	}
	if !found {
		return fmt.Errorf("%d not found", productID)
	}
	order.Items = newItems
	recalc(order)
	return nil
}

func UpdateQty(order *models.Order, productID int, qty int) error {
	if order.Status != "DRAFT" {
		return errors.New("Order with ID not draft")
	}
	if qty <= 0 {
		return errors.New("Quantity must be greater than zero")
	}
	for i := range order.Items {
		if order.Items[i].ProductID == productID {
			order.Items[i].Qty = qty
			recalc(order)
			return nil
		}
	}
	return fmt.Errorf("%d not in order", productID)
}

func SetDiscount(order *models.Order, discount int) error {
	if order.Status != "DRAFT" {
		return fmt.Errorf("Order with ID %d already closed", order.ID)
	}
	if discount <= 0 {
		return fmt.Errorf("Quantity must be greater than zero")
	}
	order.Discount = discount
	recalc(order)
	if order.Discount < order.Subtotal {
		return errors.New("Discount must be greater than or equal to subtotal")
	}
	return nil
}

func Checkout(order *models.Order, method string, paid int) error {
	if order.Status != "DRAFT" {
		return fmt.Errorf("Order with ID %d already closed", order.ID)
	}
	if order.Total < 0 {
		return fmt.Errorf("Quantity must be greater than zero")
	}
	switch method {
	case "cash":
		if paid < order.Total {
			return fmt.Errorf("Quantity must be greater than or equal to paid")
		}
		order.Paid = paid
		order.Change = paid - order.Total
		order.Method = "cash"
	case "ewallet":
		order.Paid = order.Total
		order.Change = 0
		order.Method = "ewallet"
	default:
		return fmt.Errorf("Unsupported method: %s", method)
	}
	order.Status = "PAID"
	return nil
}

func recalc(order *models.Order) {
	sub := 0
	for _, it := range order.Items {
		sub += it.Price * it.Qty
	}
	order.Subtotal = sub
	if order.Discount < 0 {
		order.Discount = 0
	}
	if order.Discount > order.Subtotal {
		order.Discount = order.Subtotal
	}
	order.Total = order.Subtotal - order.Discount
}
