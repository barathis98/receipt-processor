package utils

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParseFields(c *gin.Context) (map[string]interface{}, error) {
	rawData, err := c.GetRawData()
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(rawData, &data); err != nil {
		return nil, fmt.Errorf("invalid JSON format: %w", err)
	}

	if totalStr, ok := data["total"].(string); ok {
		total, err := strconv.ParseFloat(totalStr, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid total format: %w", err)
		}
		data["total"] = total
	} else {
		return nil, fmt.Errorf("total is required")
	}

	if items, ok := data["items"].([]interface{}); ok {
		for i, item := range items {
			itemMap, ok := item.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("invalid item format")
			}
			if priceStr, ok := itemMap["price"].(string); ok {
				price, err := strconv.ParseFloat(priceStr, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid price format for item %d: %w", i+1, err)
				}
				itemMap["price"] = price
			} else {
				return nil, fmt.Errorf("price is required for item %d", i+1)
			}
		}
		data["items"] = items
	} else {
		return nil, fmt.Errorf("items is required")
	}

	return data, nil
}
