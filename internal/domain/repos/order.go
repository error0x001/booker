package repos

import "booker/internal/domain/entites"

type OrderRepository interface {
	SaveOrder(order entites.Order) error
}
