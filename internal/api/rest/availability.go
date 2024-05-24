package rest

import (
	"booker/internal/api"
	"booker/internal/domain/entites"
	"booker/internal/infrastructure/logging"
	"booker/internal/service"
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

type RequestSearch struct {
	HotelID string    `json:"hotel_id"`
	RoomID  string    `json:"room_id"`
	From    time.Time `json:"from"`
	To      time.Time `json:"to"`
}

type ResponseSearch struct {
	HotelID string    `json:"hotel_id"`
	RoomID  string    `json:"room_id"`
	Date    time.Time `json:"date"`
	Quota   int       `json:"quota"`
}

func (h *Handler) GetAvailability(w http.ResponseWriter, r *http.Request) { //nolint:revive
	if r.Method != http.MethodGet {
		JSONError(w, api.MsgMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}
	query := r.URL.Query()
	req, err := validate(query)
	if err != nil {
		JSONError(w, err.Error(), http.StatusBadRequest)
		logging.LogErrorf("GetAvailability validation failed: %v", err)
		return
	}
	res := h.search.GetAvailability(service.Query{
		HotelID: req.HotelID,
		RoomID:  req.RoomID,
		From:    req.From,
		To:      req.To,
	})
	w.Header().Set("Content-Type", ApplicationJSONContentType)
	err = json.NewEncoder(w).Encode(convertServiceResponseToSearchQuery(res))
	if err != nil {
		logging.LogErrorf("Failed to encode response: %v", err)
	}
	logging.LogInfof("Search successfully done: %v", res)
}

func validate(query url.Values) (*RequestSearch, error) {
	hotelID := query.Get("hotel_id")
	if err := api.ValidateHotelID(hotelID); err != nil {
		return nil, err
	}

	roomID := query.Get("room_id")
	if err := api.ValidateRoomID(roomID); err != nil {
		return nil, err
	}

	from, err := time.Parse(time.RFC3339, query.Get("from"))
	if err != nil {
		return nil, api.ErrInvalidDateFrom
	}
	if err = api.ValidateCheckin(from); err != nil {
		return nil, err
	}

	to, err := time.Parse(time.RFC3339, query.Get("to"))
	if err != nil {
		return nil, api.ErrInvalidDateTo
	}

	if err = api.ValidateCheckout(from, to); err != nil {
		return nil, err
	}
	return &RequestSearch{
		HotelID: hotelID,
		RoomID:  roomID,
		From:    from,
		To:      to,
	}, nil
}

func convertServiceResponseToSearchQuery(availability []entites.RoomAvailability) []ResponseSearch {
	result := make([]ResponseSearch, len(availability))
	for i, room := range availability {
		result[i] = ResponseSearch(room)
	}
	return result
}
