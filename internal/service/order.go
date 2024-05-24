package service

import (
	"booker/internal/domain/entites"
	"booker/internal/domain/repos"
	"booker/internal/infrastructure/logging"
	"booker/internal/utils"
	"errors"
	"time"
)

var (
	ErrRoomNotAvailable = errors.New("room not available for selected dates")
	ErrSaveOrderFailed  = errors.New("could not save order")
)

type OrderService struct {
	orderRepo        repos.OrderRepository
	availabilityRepo repos.AvailabilityRepository
}

func NewOrderService(orderRepo repos.OrderRepository, availabilityRepo repos.AvailabilityRepository) *OrderService {
	return &OrderService{
		orderRepo:        orderRepo,
		availabilityRepo: availabilityRepo,
	}
}

func (s *OrderService) Create(order entites.Order) error {
	daysToBook := utils.DaysBetween(order.From, order.To)

	s.availabilityRepo.Lock()
	defer s.availabilityRepo.Unlock()

	oldQuotas := make(map[string]int)

	for _, day := range daysToBook {
		availability, err := s.availabilityRepo.GetAvailability(order.HotelID, order.RoomID, day)
		if err != nil {
			s.rollbackAvailabilities(daysToBook, oldQuotas, order)
			return err
		}
		if availability.Quota < 1 {
			s.rollbackAvailabilities(daysToBook, oldQuotas, order)
			return ErrRoomNotAvailable
		}

		key := utils.MakeKey(order.HotelID, order.RoomID, day)
		oldQuotas[key] = availability.Quota

		err = s.availabilityRepo.UpdateAvailability(order.HotelID, order.RoomID, day, availability.Quota-1)
		if err != nil {
			s.rollbackAvailabilities(daysToBook, oldQuotas, order)
			return err
		}
	}

	err := s.orderRepo.SaveOrder(order)
	if err != nil {
		s.rollbackAvailabilities(daysToBook, oldQuotas, order)
		return ErrSaveOrderFailed
	}

	return nil
}

func (s *OrderService) rollbackAvailabilities(daysToBook []time.Time, oldQuotas map[string]int, order entites.Order) {
	for _, day := range daysToBook {
		key := utils.MakeKey(order.HotelID, order.RoomID, day)
		if oldQuota, ok := oldQuotas[key]; ok {
			err := s.availabilityRepo.UpdateAvailability(order.HotelID, order.RoomID, day, oldQuota)
			if err != nil {
				logging.LogErrorf("Failed to rollback availability for key %s: %v\n", key, err)
			}
		} else {
			logging.LogInfof("No old quota found for key %s during rollback.\n", key)
		}
	}
}
