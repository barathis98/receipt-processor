package store

import (
	"fmt"
	"sync"
)

type PointsStore struct {
	data sync.Map
}

func NewPointsStore() *PointsStore {
	return &PointsStore{}
}

func (s *PointsStore) Add(id string, points int) error {
	if id == "" {
		return fmt.Errorf("ID cannot be empty")
	}

	if points < 0 {
		return fmt.Errorf("points cannot be negative")
	}

	s.data.Store(id, points)
	return nil
}

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

func (s *PointsStore) Update(id string, points int) error {
	if id == "" {
		return fmt.Errorf("ID cannot be empty")
	}

	if points < 0 {
		return fmt.Errorf("points cannot be negative")
	}

	_, exists := s.data.Load(id)
	if !exists {
		return fmt.Errorf("points for receipt ID %s not found, cannot update", id)
	}

	s.data.Store(id, points)
	return nil
}

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

	if len(pointsMap) == 0 {
		return nil, fmt.Errorf("no points found in the store")
	}

	return pointsMap, nil
}
