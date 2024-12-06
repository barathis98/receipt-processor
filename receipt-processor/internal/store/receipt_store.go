package store

import (
	"fmt"
	"receipt-processor/internal/model"
	"sync"
)

type ReceiptStore struct {
	data sync.Map
}

func NewReceiptStore() *ReceiptStore {
	return &ReceiptStore{}
}

func (s *ReceiptStore) Add(receipt *model.Receipt) error {
	if receipt == nil {
		return fmt.Errorf("receipt cannot be nil")
	}
	s.data.Store(receipt.ID, receipt)
	return nil
}

func (s *ReceiptStore) Get(id string) (*model.Receipt, error) {
	if id == "" {
		return nil, fmt.Errorf("ID cannot be empty")
	}

	value, exists := s.data.Load(id)
	if !exists {
		return nil, fmt.Errorf("receipt with ID %s not found", id)
	}

	receipt, ok := value.(*model.Receipt)
	if !ok {
		return nil, fmt.Errorf("invalid receipt data for ID %s", id)
	}

	return receipt, nil
}

func (s *ReceiptStore) Delete(id string) error {
	if id == "" {
		return fmt.Errorf("ID cannot be empty")
	}

	_, exists := s.data.Load(id)
	if !exists {
		return fmt.Errorf("receipt with ID %s not found", id)
	}

	s.data.Delete(id)
	return nil
}

// List returns all receipts from the store with error handling.
func (s *ReceiptStore) List() ([]*model.Receipt, error) {
	var receipts []*model.Receipt
	s.data.Range(func(key, value interface{}) bool {
		receipt, ok := value.(*model.Receipt)
		if !ok {
			// If the data type is incorrect, return an error
			return false
		}
		receipts = append(receipts, receipt)
		return true
	})

	// Check if any receipts were found
	if len(receipts) == 0 {
		return nil, fmt.Errorf("no receipts found")
	}

	return receipts, nil
}
