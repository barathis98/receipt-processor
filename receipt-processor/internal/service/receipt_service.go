package service

import (
	"fmt"
	"math"
	"receipt-processor/internal/model"
	"receipt-processor/internal/store"
	"receipt-processor/internal/utils"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

func SaveReceipt(payload model.Receipt) (string, error) {
	newReceipt := model.NewReceipt(
		payload.Retailer,
		payload.PurchaseDate,
		payload.PurchaseTime,
		payload.Items,
		payload.Total,
	)

	if err := store.RS.Add(newReceipt); err != nil {
		utils.Logger.Error("Failed to add receipt to store", zap.String("receipt_id", newReceipt.ID), zap.Error(err))
		return "", fmt.Errorf("could not add receipt to store: %w", err)
	}

	points, err := CalculatePoints(*newReceipt)
	if err != nil {
		utils.Logger.Error("Failed to calculate points", zap.String("receipt_id", newReceipt.ID), zap.Error(err))
		return "", fmt.Errorf("could not calculate points for receipt: %w", err)
	}

	if err := store.PS.Add(newReceipt.ID, points); err != nil {
		utils.Logger.Error("Failed to add points to store", zap.String("receipt_id", newReceipt.ID), zap.Error(err))
		return "", fmt.Errorf("could not add points to store for receipt: %w", err)
	}

	utils.Logger.Info("Receipt saved successfully", zap.String("receipt_id", newReceipt.ID))
	return newReceipt.ID, nil
}

func GetReceiptById(id string) (model.Receipt, error) {
	receipt, exists := store.RS.Get(id)
	if exists != nil {
		return model.Receipt{}, fmt.Errorf("receipt with ID %s not found", id)
	}
	return *receipt, nil
}

func GetPoints(id string) (int, error) {
	points, exists := store.PS.Get(id)
	if exists != nil {
		return 0, fmt.Errorf("points for receipt with ID %s not found", id)
	}
	return points, nil
}

func CalculatePoints(receipt model.Receipt) (int, error) {
	points := 0

	for _, ch := range receipt.Retailer {
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') {
			points++
		}
	}

	totalFloat, err := strconv.ParseFloat(receipt.Total, 64)
	if err != nil {
		fmt.Println("Error converting total to float:", err)
	}

	if totalFloat <= 0 {
		return 0, fmt.Errorf("invalid total amount: must be positive")
	}

	if totalFloat == float64(int(totalFloat)) {
		points += 50
	}

	if math.Mod(totalFloat, 0.25) == 0 {
		points += 25
	}

	if len(receipt.Items) > 0 {
		points += (len(receipt.Items) / 2) * 5
	}

	purchaseDate, err := time.Parse("2006-01-02", receipt.PurchaseDate)
	if err != nil {
		return points, fmt.Errorf("invalid purchase date format: %v", err)
	}

	if purchaseDate.Day()%2 != 0 {
		points += 6
	}

	for _, item := range receipt.Items {
		description := strings.TrimSpace(item.ShortDesc)
		if len(description)%3 == 0 {

			itemPrice, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				fmt.Printf("Error parsing price for item %s: %v\n", description, err)
				continue
			}

			points += int(math.Ceil(itemPrice * 0.2))
		}
	}

	purchaseTime, err := time.Parse("15:04", receipt.PurchaseTime)
	if err != nil {
		return points, fmt.Errorf("invalid purchase time format: %v", err)
	}

	if (purchaseTime.Hour() == 14 && purchaseTime.Minute() > 0) || (purchaseTime.Hour() == 15) {
		points += 10
	}

	return points, nil
}
