package service

import (
	"booker/internal/domain/entites"
	"errors"
	"sync"
	"testing"
	"time"
)

type MockOrderRepository struct {
	saveOrderFunc func(order entites.Order) error
	mutex         sync.Mutex
}

func (m *MockOrderRepository) SaveOrder(order entites.Order) error {
	return m.saveOrderFunc(order)
}

type MockAvailabilityRepository struct {
	getAvailabilityFunc    func(hotelID, roomID string, date time.Time) (*entites.RoomAvailability, error)
	updateAvailabilityFunc func(hotelID, roomID string, date time.Time, quota int) error
	mutex                  sync.Mutex
}

func (m *MockAvailabilityRepository) Lock() {
	m.mutex.Lock()
}

func (m *MockAvailabilityRepository) Unlock() {
	m.mutex.Unlock()
}

func (m *MockAvailabilityRepository) GetAvailability(hotelID, roomID string, date time.Time) (*entites.RoomAvailability, error) {
	return m.getAvailabilityFunc(hotelID, roomID, date)
}

func (m *MockAvailabilityRepository) UpdateAvailability(hotelID, roomID string, date time.Time, quota int) error {
	return m.updateAvailabilityFunc(hotelID, roomID, date, quota)
}

func TestOrderService_Create(t *testing.T) {
	orderRepo := &MockOrderRepository{}
	availabilityRepo := &MockAvailabilityRepository{}
	service := NewOrderService(orderRepo, availabilityRepo)

	order := entites.Order{
		HotelID:   "reddison",
		RoomID:    "lux",
		UserEmail: "user@2gis.ru",
		From:      time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
		To:        time.Date(2023, 10, 3, 0, 0, 0, 0, time.UTC),
	}

	t.Run("ok", func(t *testing.T) {
		availabilityRepo.getAvailabilityFunc = func(hotelID, roomID string, date time.Time) (*entites.RoomAvailability, error) {
			return &entites.RoomAvailability{
				HotelID: hotelID,
				RoomID:  roomID,
				Date:    date,
				Quota:   1,
			}, nil
		}

		availabilityRepo.updateAvailabilityFunc = func(hotelID, roomID string, date time.Time, quota int) error {
			return nil
		}

		orderRepo.saveOrderFunc = func(o entites.Order) error {
			return nil
		}

		err := service.Create(order)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("room_not_available", func(t *testing.T) {
		availabilityRepo.getAvailabilityFunc = func(hotelID, roomID string, date time.Time) (*entites.RoomAvailability, error) {
			return &entites.RoomAvailability{
				HotelID: hotelID,
				RoomID:  roomID,
				Date:    date,
				Quota:   0,
			}, nil
		}

		err := service.Create(order)
		if !errors.Is(err, ErrRoomNotAvailable) {
			t.Errorf("expected error %v, got %v", ErrRoomNotAvailable, err)
		}
	})

	t.Run("saving_failed", func(t *testing.T) {
		availabilityRepo.getAvailabilityFunc = func(hotelID, roomID string, date time.Time) (*entites.RoomAvailability, error) {
			return &entites.RoomAvailability{
				HotelID: hotelID,
				RoomID:  roomID,
				Date:    date,
				Quota:   1,
			}, nil
		}

		availabilityRepo.updateAvailabilityFunc = func(hotelID, roomID string, date time.Time, quota int) error {
			return nil
		}

		orderRepo.saveOrderFunc = func(o entites.Order) error {
			return errors.New("save error")
		}

		err := service.Create(order)
		if !errors.Is(err, ErrSaveOrderFailed) {
			t.Errorf("expected error %v, got %v", ErrSaveOrderFailed, err)
		}
	})
}

func TestOrderService_ConcurrentCreate(t *testing.T) {
	orderRepo := &MockOrderRepository{}
	availabilityRepo := &MockAvailabilityRepository{}
	service := NewOrderService(orderRepo, availabilityRepo)

	order1 := entites.Order{
		HotelID:   "reddison",
		RoomID:    "lux",
		UserEmail: "user1@2gis.ru",
		From:      time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
		To:        time.Date(2023, 10, 3, 0, 0, 0, 0, time.UTC),
	}

	order2 := entites.Order{
		HotelID:   "reddison",
		RoomID:    "lux",
		UserEmail: "user2@2gis.ru",
		From:      time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
		To:        time.Date(2023, 10, 3, 0, 0, 0, 0, time.UTC),
	}

	availabilityCounter := 2

	availabilityRepo.getAvailabilityFunc = func(hotelID, roomID string, date time.Time) (*entites.RoomAvailability, error) {
		return &entites.RoomAvailability{
			HotelID: hotelID,
			RoomID:  roomID,
			Date:    date,
			Quota:   availabilityCounter,
		}, nil
	}

	availabilityRepo.updateAvailabilityFunc = func(hotelID, roomID string, date time.Time, quota int) error {
		availabilityCounter = quota
		return nil
	}

	orderRepo.saveOrderFunc = func(o entites.Order) error {
		return nil
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := service.Create(order1)
		if err != nil && !errors.Is(err, ErrRoomNotAvailable) {
			t.Errorf("expected no error or room not available, got %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		err := service.Create(order2)
		if err != nil && !errors.Is(err, ErrRoomNotAvailable) {
			t.Errorf("expected no error or room not available, got %v", err)
		}
	}()

	wg.Wait()
}
