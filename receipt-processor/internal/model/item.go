package model

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Item struct {
	ShortDesc string  `json:"shortDescription"`
	Price     float64 `json:"price"`
}

func (i *Item) UnmarshalJSON(data []byte) error {

	type Alias Item
	aux := &struct {
		Price interface{} `json:"price"`
		*Alias
	}{
		Alias: (*Alias)(i),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	switch v := aux.Price.(type) {
	case string:
		parsedPrice, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return fmt.Errorf("invalid price value: %v", v)
		}
		i.Price = parsedPrice
	case float64:
		i.Price = v
	default:
		return fmt.Errorf("price should be a string or float64, got %T", v)
	}

	return nil
}
