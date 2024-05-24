package utils

import (
	"errors"
	"fmt"
	"time"
)

func MakeKey(hotelID, roomID string, date time.Time) string {
	return fmt.Sprintf("%s_%s_%s", hotelID, roomID, date.Format(time.DateOnly))
}

func ParseKey(key string) (hotelID, roomID string, date time.Time, err error) {
	var dateString string
	n, err := fmt.Sscanf(key, "%s_%s_%s", &hotelID, &roomID, &dateString)
	if err != nil || n != 3 {
		return "", "", time.Time{}, errors.New("invalid key format")
	}
	date, err = time.Parse(time.DateOnly, dateString)
	if err != nil {
		return "", "", time.Time{}, err
	}
	return hotelID, roomID, date, nil
}
