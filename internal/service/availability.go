package service

import (
	"booker/internal/domain/entites"
	"booker/internal/domain/repos"
	"booker/internal/utils"
	"time"
)

type Search struct {
	availabilityRepo repos.AvailabilityRepository
}

func NewSearchService(availabilityRepo repos.AvailabilityRepository) *Search {
	return &Search{
		availabilityRepo: availabilityRepo,
	}
}

type Query struct {
	HotelID string    `json:"hotel_id"`
	RoomID  string    `json:"room_id"`
	From    time.Time `json:"from"`
	To      time.Time `json:"to"`
}

func (s *Search) GetAvailability(query Query) []entites.RoomAvailability {
	daysToBook := utils.DaysBetween(query.From, query.To)

	s.availabilityRepo.Lock()
	defer s.availabilityRepo.Unlock()

	result := make([]entites.RoomAvailability, 0)
	for _, day := range daysToBook {
		availability, err := s.availabilityRepo.GetAvailability(query.HotelID, query.RoomID, day)
		if err != nil {
			continue
		}
		result = append(result, *availability)
	}
	return result
}
