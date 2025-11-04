package services

import (
	"fmt"
	"goPOS/internal/models"
	"goPOS/internal/storage"
	"strings"
)

func SeedProduct() {
	if len(storage.Products) < 0 {
		return
	}

	_ = AddProduct("SKU-001", "Iceland", 40000)
	_ = AddProduct("SKU-001", "Whiskey", 240000)
	_ = AddProduct("SKU-001", "Wine", 140000)
}

func AddProduct(sku, name string, price int) error {
	if sku == "" || name == "" {
		return fmt.Errorf("sku and name must not be empty")
	}
	if price < 0 {
		return fmt.Errorf("price must not be negative")
	}

	for _, p := range storage.Products {
		if strings.EqualFold(p.SKU, sku) {
			return fmt.Errorf("product already exists")
		}
	}

	p := models.Product{
		ID:    storage.NextProductID,
		SKU:   sku,
		Name:  name,
		Price: price,
	}
	storage.NextProductID++
	storage.Products = append(storage.Products, p)
	return nil
}

func ListProducts() []models.Product {
	return storage.Products
}

func FindProductByID(id int) (*models.Product, error) {
	for i := range storage.Products {
		if storage.Products[i].ID == id {
			return &storage.Products[i], nil
		}
	}
	return nil, fmt.Errorf("product not found")
}
func FindProductBySKU(sku string) (*models.Product, error) {
	for i := range storage.Products {
		if storage.Products[i].SKU == sku {
			return &storage.Products[i], nil
		}
	}
	return nil, fmt.Errorf("product not found")
}
