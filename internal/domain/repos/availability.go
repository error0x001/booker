package repos

import (
	"booker/internal/domain/entites"
	"time"
)

type AvailabilityRepository interface {
	GetAvailability(hotelID, roomID string, date time.Time) (*entites.RoomAvailability, error)
	UpdateAvailability(hotelID, roomID string, date time.Time, quota int) error
	Lock()
	Unlock()
}
