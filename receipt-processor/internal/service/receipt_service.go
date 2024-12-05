package service

import (
	"fmt"
	"math"
	"receipt-processor/internal/model"
	"receipt-processor/internal/store"
	"receipt-processor/internal/utils"
	"strings"
	"time"

	"go.uber.org/zap"
)

func SaveReceipt(payload model.Receipt) (string, error) {
	// Create a new receipt using the payload
	newReceipt := model.NewReceipt(
		payload.Retailer,
		payload.PurchaseDate,
		payload.PurchaseTime,
		payload.Items,
		payload.Total,
	)

	// Try adding the new receipt to the ReceiptStore
	if err := store.RS.Add(newReceipt); err != nil {
		// Log the error
		utils.Logger.Error("Failed to add receipt to store", zap.String("receipt_id", newReceipt.ID), zap.Error(err))
		return "", fmt.Errorf("could not add receipt to store: %w", err)
	}

	// Calculate points for the receipt
	points, err := CalculatePoints(*newReceipt)
	if err != nil {
		// Log the error
		utils.Logger.Error("Failed to calculate points", zap.String("receipt_id", newReceipt.ID), zap.Error(err))
		return "", fmt.Errorf("could not calculate points for receipt: %w", err)
	}

	// Try adding the points to the PointsStore
	if err := store.PS.Add(newReceipt.ID, points); err != nil {
		// Log the error
		utils.Logger.Error("Failed to add points to store", zap.String("receipt_id", newReceipt.ID), zap.Error(err))
		return "", fmt.Errorf("could not add points to store for receipt: %w", err)
	}

	// Return the receipt ID if successful
	utils.Logger.Info("Receipt saved successfully", zap.String("receipt_id", newReceipt.ID))
	return newReceipt.ID, nil
}

func GetReceiptById(id string) (model.Receipt, error) {
	// Retrieve the receipt from the store using the ID
	receipt, exists := store.RS.Get(id)
	if exists != nil {
		// Return an error if the receipt is not found
		return model.Receipt{}, fmt.Errorf("receipt with ID %s not found", id)
	}
	// Return the found receipt
	return *receipt, nil
}

func GetPoints(id string) (int, error) {
	// Retrieve the points from the store using the receipt ID
	points, exists := store.PS.Get(id)
	if exists != nil {
		// Return an error if the points for the receipt ID are not found
		return 0, fmt.Errorf("points for receipt with ID %s not found", id)
	}
	// Return the points if found
	return points, nil
}

// CalculatePoints calculates points based on the receipt details
func CalculatePoints(receipt model.Receipt) (int, error) {
	points := 0

	// Validate Retailer Name
	for _, ch := range receipt.Retailer {
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') {
			points++
		}
	}

	// Validate the Total
	if receipt.Total <= 0 {
		return 0, fmt.Errorf("invalid total amount: must be positive")
	}

	// If the total is an integer, add 50 points
	if receipt.Total == float64(int(receipt.Total)) {
		points += 50
	}

	// If the total is divisible by 0.25, add 25 points
	if math.Mod(receipt.Total, 0.25) == 0 {
		points += 25
	}

	// Add points for items
	if len(receipt.Items) > 0 {
		points += (len(receipt.Items) / 2) * 5
	}

	// Validate Purchase Date
	purchaseDate, err := time.Parse("2006-01-02", receipt.PurchaseDate)
	if err != nil {
		return points, fmt.Errorf("invalid purchase date format: %v", err)
	}

	// If the day of the month is odd, add 6 points
	if purchaseDate.Day()%2 != 0 {
		points += 6
	}

	// Calculate points for each item based on description length
	for _, item := range receipt.Items {
		description := strings.TrimSpace(item.ShortDesc)
		if len(description)%3 == 0 {
			points += int(math.Ceil(item.Price * 0.2))
		}
	}

	// Validate Purchase Time
	purchaseTime, err := time.Parse("15:04", receipt.PurchaseTime)
	if err != nil {
		return points, fmt.Errorf("invalid purchase time format: %v", err)
	}

	// Add 10 points for purchases between 2:00 PM to 3:00 PM
	if (purchaseTime.Hour() == 14 && purchaseTime.Minute() > 0) || (purchaseTime.Hour() == 15) {
		points += 10
	}

	return points, nil
}
