package storage

import (
	"booker/internal/domain/entites"
	"booker/internal/utils"
	"errors"
	"sync"
	"time"
)

var errAvailabilityNotFound = errors.New("availability not found")

type MemoryAvailabilityRepository struct {
	mu           sync.RWMutex
	Availability map[string]entites.RoomAvailability
}

func NewMemoryAvailabilityRepository() *MemoryAvailabilityRepository {
	return &MemoryAvailabilityRepository{
		mu:           sync.RWMutex{},
		Availability: make(map[string]entites.RoomAvailability),
	}
}

func (repo *MemoryAvailabilityRepository) GetAvailability(hotelID, roomID string, date time.Time) (*entites.RoomAvailability, error) {
	key := utils.MakeKey(hotelID, roomID, date)
	if availability, ok := repo.Availability[key]; ok {
		return &availability, nil
	}
	return nil, errAvailabilityNotFound
}

func (repo *MemoryAvailabilityRepository) UpdateAvailability(hotelID, roomID string, date time.Time, quota int) error {
	key := utils.MakeKey(hotelID, roomID, date)
	if availability, ok := repo.Availability[key]; ok {
		availability.Quota = quota
		repo.Availability[key] = availability
		return nil
	}
	return errAvailabilityNotFound
}

func (repo *MemoryAvailabilityRepository) Lock() {
	repo.mu.Lock()
}

func (repo *MemoryAvailabilityRepository) Unlock() {
	repo.mu.Unlock()
}
