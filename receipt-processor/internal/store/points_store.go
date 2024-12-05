package store

import (
	"fmt"
	"sync"
)

type PointsStore struct {
	data sync.Map // Thread-safe map for storing receipt points
}

// Initialize a new PointsStore
func NewPointsStore() *PointsStore {
	return &PointsStore{}
}

// Add points for a receipt ID with error handling
func (s *PointsStore) Add(id string, points int) error {
	if id == "" {
		return fmt.Errorf("ID cannot be empty")
	}

	// Ensure points is not negative (optional validation)
	if points < 0 {
		return fmt.Errorf("points cannot be negative")
	}

	s.data.Store(id, points)
	return nil
}

// Get points for a receipt ID with error handling
func (s *PointsStore) Get(id string) (int, error) {
	if id == "" {
		return 0, fmt.Errorf("ID cannot be empty")
	}

	value, exists := s.data.Load(id)
	if !exists {
		return 0, fmt.Errorf("points for receipt ID %s not found", id)
	}

	points, ok := value.(int)
	if !ok {
		return 0, fmt.Errorf("invalid points data for receipt ID %s", id)
	}

	return points, nil
}

// Update points for a receipt ID with error handling
func (s *PointsStore) Update(id string, points int) error {
	if id == "" {
		return fmt.Errorf("ID cannot be empty")
	}

	// Ensure points is not negative (optional validation)
	if points < 0 {
		return fmt.Errorf("points cannot be negative")
	}

	// Check if the points exist before updating (optional)
	_, exists := s.data.Load(id)
	if !exists {
		return fmt.Errorf("points for receipt ID %s not found, cannot update", id)
	}

	s.data.Store(id, points)
	return nil
}

// Delete points for a receipt ID with error handling
func (s *PointsStore) Delete(id string) error {
	if id == "" {
		return fmt.Errorf("ID cannot be empty")
	}

	_, exists := s.data.Load(id)
	if !exists {
		return fmt.Errorf("points for receipt ID %s not found, cannot delete", id)
	}

	s.data.Delete(id)
	return nil
}

// List all receipt IDs and points with error handling
func (s *PointsStore) List() (map[string]int, error) {
	pointsMap := make(map[string]int)
	s.data.Range(func(key, value interface{}) bool {
		id, idOk := key.(string)
		points, pointsOk := value.(int)
		if idOk && pointsOk {
			pointsMap[id] = points
		}
		return true
	})

	// If the points map is empty, return an error
	if len(pointsMap) == 0 {
		return nil, fmt.Errorf("no points found in the store")
	}

	return pointsMap, nil
}
