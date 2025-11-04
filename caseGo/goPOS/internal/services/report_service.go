package services

import (
	"goPOS/internal/models"
	"goPOS/internal/storage"
)

func GenerateReport() models.Report {
	var r models.Report
	for _, o := range storage.Orders {
		if o.Status == "PAID" {
			r.Transaction++
			r.Revenue += o.Total
		}
	}
	return r
}
