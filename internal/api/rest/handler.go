package rest

import (
	"booker/internal/service"
)

type Handler struct {
	search *service.Search
	order  *service.OrderService
}

func NewHandler(orderService *service.OrderService, searchService *service.Search) *Handler {
	return &Handler{order: orderService, search: searchService}
}
