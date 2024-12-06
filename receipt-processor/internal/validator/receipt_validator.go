package validator

import (
	"errors"
	"fmt"
	"receipt-processor/internal/model"
	"strconv"
	"time"
)

func ValidateReceipt(receipt model.Receipt) error {
	if receipt.Retailer == "" {
		return errors.New("retailer is required")
	}
	if receipt.PurchaseDate == "" {
		return errors.New("purchaseDate is required")
	}
	if receipt.PurchaseTime == "" {
		return errors.New("purchaseTime is required")
	}
	if len(receipt.Items) == 0 {
		return errors.New("at least one item is required")
	}

	if _, err := time.Parse("2006-01-02", receipt.PurchaseDate); err != nil {
		return errors.New("invalid purchaseDate format, expected YYYY-MM-DD")
	}

	if _, err := time.Parse("15:04", receipt.PurchaseTime); err != nil {
		return errors.New("invalid purchaseTime format, expected HH:MM")
	}

	total, err := strconv.ParseFloat(receipt.Total, 64)
	if err != nil {
		return errors.New("invalid total: must be a numeric value")
	}
	if total <= 0 {
		return errors.New("total must be greater than zero")
	}

	for i, item := range receipt.Items {
		if item.ShortDesc == "" {
			return fmt.Errorf("shortDescription is required for item %d", i+1)
		}

		price, err := strconv.ParseFloat(item.Price, 64)
		if err != nil {
			return fmt.Errorf("invalid price for item %d: must be a numeric value", i+1)
		}

		if price <= 0 {
			return fmt.Errorf("price must be greater than zero for item %d", i+1)
		}
	}

	return nil
}
