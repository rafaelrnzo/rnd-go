package storage

import "goPOS/internal/models"

var Products []models.Product
var Orders []models.Order

var NextProductID = 1
var NextOrderID = 1

func ResetSession() {
	Products = nil
	Orders = nil
	NextProductID = 1
	NextOrderID = 1
}
