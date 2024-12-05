package service

import (
	"fmt"
	"math"
	"receipt-processor/internal/model"
	"strings"
	"time"

	"github.com/google/uuid"
)

var receipts = make(map[string]model.Receipt)

func SaveReceipt(receipt model.Receipt) (string, error) {
	receipt.ID = uuid.New().String()
	receipt.Points = CalculatePoints(receipt)
	receipts[receipt.ID] = receipt
	return receipt.ID, nil
}

func GetReceiptById(id string) (model.Receipt, error) {
	receipt, exists := receipts[id]
	if !exists {
		return model.Receipt{}, fmt.Errorf("receipt not found")
	}
	return receipt, nil
}

func GetPoints(id string) (int, error) {
	receipt, exists := receipts[id]
	if !exists {
		return 0, fmt.Errorf("receipt not found")
	}
	return receipt.Points, nil
}

func CalculatePoints(receipt model.Receipt) int {
	points := 0
	for _, ch := range receipt.Retailer {
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') {
			points++
		}
	}

	if receipt.Total == float64(int(receipt.Total)) {
		points += 50
	}

	if math.Mod(receipt.Total, 0.25) == 0 {
		points += 25
	}

	points += (len(receipt.Items) / 2) * 5

	purchaseDate, err := time.Parse("2006-01-02", receipt.PurchaseDate)
	if err == nil && purchaseDate.Day()%2 != 0 {
		points += 6
	}

	for _, item := range receipt.Items {
		description := strings.TrimSpace(item.ShortDesc)
		if len(description)%3 == 0 {
			points += int(math.Ceil(item.Price * 0.2))
		}
	}

	purchaseTime, err := time.Parse("15:04", receipt.PurchaseTime)
	if err == nil && ((purchaseTime.Hour() == 14 && purchaseTime.Minute() > 0) ||
		(purchaseTime.Hour() == 15)) {
		points += 10
	}

	return points
}
