package model

import (
	"github.com/google/uuid"
)

type Receipt struct {
	ID           string  `json:"id,omitempty"`
	Retailer     string  `json:"retailer" binding:"required"`
	PurchaseDate string  `json:"purchaseDate" binding:"required"`
	PurchaseTime string  `json:"purchaseTime" binding:"required"`
	Items        []Item  `json:"items" binding:"required"`
	Total        float64 `json:"total" binding:"required"`
}

func NewReceipt(retailer string, purchaseDate string, purchaseTime string, items []Item, total float64) *Receipt {
	return &Receipt{
		ID:           uuid.New().String(),
		Retailer:     retailer,
		PurchaseDate: purchaseDate,
		PurchaseTime: purchaseTime,
		Items:        items,
		Total:        total,
	}
}
