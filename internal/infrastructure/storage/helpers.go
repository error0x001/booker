package storage

import (
	"booker/internal/config"
	"booker/internal/domain/entites"
	"booker/internal/domain/repos"
	"booker/internal/utils"
	"fmt"
	"time"
)

func prepareMemoryAvailability(conf config.Config, availabilityRepo *MemoryAvailabilityRepository) {
	if !conf.IsPrepareAvailabilityRequired {
		return
	}

	currentYear := time.Now().Year()
	hotelID := "reddison"
	roomID := "lux"

	for month := 1; month <= 12; month++ {
		for day := 1; day <= 31; day++ {
			date := utils.Date(currentYear, month, day)
			if date.Month() != time.Month(month) {
				continue
			}

			key := utils.MakeKey(hotelID, roomID, date)
			availabilityRepo.Availability[key] = entites.RoomAvailability{
				HotelID: hotelID,
				RoomID:  roomID,
				Date:    date,
				Quota:   utils.RandRange(0, 10),
			}
		}
	}
}

func CreateRepositories(conf config.Config) (repos.OrderRepository, repos.AvailabilityRepository, error) {
	switch conf.StorageType {
	case config.Memory:
		orderRepo := NewMemoryOrderRepository()
		availabilityRepo := NewMemoryAvailabilityRepository()
		prepareMemoryAvailability(conf, availabilityRepo)
		return orderRepo, availabilityRepo, nil
	case config.External:
		// TODO add connections to external databases
		return nil, nil, fmt.Errorf("external storage not implemented yet")
	default:
		return nil, nil, fmt.Errorf("unknown storage type: %s", conf.StorageType)
	}
}
