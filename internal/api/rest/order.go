package rest

import (
	"booker/internal/api"
	"booker/internal/domain/entites"
	"booker/internal/infrastructure/logging"
	"booker/internal/service"
	"booker/internal/utils"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type RequestOrder struct {
	HotelID   string    `json:"hotel_id"`
	RoomID    string    `json:"room_id"`
	UserEmail string    `json:"email"`
	From      time.Time `json:"from"`
	To        time.Time `json:"to"`
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		JSONError(w, api.MsgMethodNotAllowed, http.StatusMethodNotAllowed)
	}
	var req RequestOrder
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		JSONError(w, api.MsgInvalidJSONFormat, http.StatusBadRequest)
		logging.LogErrorf("JSON decode failed: %v", err)
		return
	}

	err = req.validate()
	if err != nil {
		JSONError(w, err.Error(), http.StatusBadRequest)
		logging.LogErrorf("RequestOrder validation failed: %v", err)
		return
	}

	err = h.order.Create(convertRequestToOrder(req))
	if err != nil {
		message, statusCode := convertServiceErrorToExternal(err)
		JSONError(w, message, statusCode)
		logging.LogErrorf("Order creation failed: %v", err)
		return
	}

	w.Header().Set("Content-Type", ApplicationJSONContentType)
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(req)
	if err != nil {
		logging.LogErrorf("Failed to encode response: %v", err)
	}
	logging.LogInfof("Order successfully created: %v", req)
}

func (r *RequestOrder) validate() error {
	if err := api.ValidateHotelID(r.HotelID); err != nil {
		return err
	}

	if err := api.ValidateRoomID(r.RoomID); err != nil {
		return err
	}
	if !utils.IsEmailValid(r.UserEmail) {
		return api.ErrInvalidUserEmail
	}

	if err := api.ValidateCheckin(r.From); err != nil {
		return err
	}

	return api.ValidateCheckout(r.From, r.To)
}

func convertRequestToOrder(order RequestOrder) entites.Order {
	return entites.Order(order)
}

func convertServiceErrorToExternal(err error) (string, int) {
	if errors.Is(err, service.ErrRoomNotAvailable) {
		return api.MsgRoomNotAvailable, http.StatusConflict
	}
	if errors.Is(err, service.ErrSaveOrderFailed) {
		return api.MsgSaveOrderFailed, http.StatusServiceUnavailable
	}
	return api.MsgUnknownError, http.StatusInternalServerError
}
