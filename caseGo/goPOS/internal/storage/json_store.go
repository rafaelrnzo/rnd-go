package storage

import (
	"encoding/json"
	"goPOS/internal/models"
	"os"
)

type snapshot struct {
	Products      []models.Product `json:"products"`
	Orders        []models.Order   `json:"orders"`
	NextProductID int              `json:"next_product_id"`
	NextOrderID   int              `json:"next_order_id"`
}

func SaveToFile(filename string) error {
	snap := snapshot{
		Products:      Products,
		Orders:        Orders,
		NextProductID: NextProductID,
		NextOrderID:   NextOrderID,
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(snap)
}

func LoadFromFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	var snap snapshot
	err = dec.Decode(&snap)
	if err != nil {
		return err
	}
	Products = snap.Products
	Orders = snap.Orders
	NextProductID = snap.NextProductID
	NextOrderID = snap.NextOrderID
	return nil
}
