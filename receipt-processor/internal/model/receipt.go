package model

import "github.com/google/uuid"

type Receipt struct {
	ID           string `json:"id,omitempty"`
	Retailer     string `json:"retailer" binding:"required"`
	PurchaseDate string `json:"purchaseDate" binding:"required"`
	PurchaseTime string `json:"purchaseTime" binding:"required"`
	Items        []Item `json:"items" binding:"required"`
	Total        string `json:"total" binding:"required"`
}

// func (t *Receipt) UnmarshalJSON(data []byte) error {
// 	type Alias Receipt
// 	aux := &struct {
// 		Total string `json:"total"`
// 		*Alias
// 	}{
// 		Alias: (*Alias)(t),
// 	}

// 	if err := json.Unmarshal(data, &aux); err != nil {
// 		return err
// 	}

// 	if aux.Total != "" {
// 		strValue := aux.Total
// 		if floatValue, err := strconv.ParseFloat(strValue, 64); err == nil {
// 			t.Total = floatValue
// 		} else {
// 			return fmt.Errorf("invalid total value")
// 		}
// 	}

// 	return nil
// }

func NewReceipt(retailer string, purchaseDate string, purchaseTime string, items []Item, total string) *Receipt {
	return &Receipt{
		ID:           uuid.New().String(),
		Retailer:     retailer,
		PurchaseDate: purchaseDate,
		PurchaseTime: purchaseTime,
		Items:        items,
		Total:        total,
	}
}
