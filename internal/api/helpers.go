package api

import (
	"errors"
	"time"
)

const (
	MsgInvalidJSONFormat = "invalid_json_format"
	MsgMethodNotAllowed  = "method_not_allowed"
	MsgRoomNotAvailable  = "room_not_available"
	MsgSaveOrderFailed   = "save_order_failed"
	MsgUnknownError      = "unknown_error"
)

var (
	errInvalidHotelID   = errors.New("invalid_hotel_id")
	errInvalidRoomID    = errors.New("invalid_room_id")
	ErrInvalidUserEmail = errors.New("invalid_user_email")
	ErrInvalidDateFrom  = errors.New("invalid_from")
	ErrInvalidDateTo    = errors.New("invalid_to")
)

func ValidateHotelID(hotelID string) error {
	if hotelID != "reddison" {
		return errInvalidHotelID
	}
	return nil
}

func ValidateRoomID(roomID string) error {
	if roomID != "lux" {
		return errInvalidRoomID
	}
	return nil
}

func ValidateCheckin(from time.Time) error {
	today := time.Now().Truncate(24 * time.Hour)
	fromDate := from.Truncate(24 * time.Hour)
	if today.Sub(fromDate) >= 0 {
		return ErrInvalidDateFrom
	}
	return nil
}

func ValidateCheckout(from, to time.Time) error {
	if from.After(to) {
		return ErrInvalidDateTo
	}
	return nil
}
