package storage

import (
	"booker/internal/domain/entites"
	"sync"
)

type MemoryOrderRepository struct {
	mu     sync.RWMutex
	orders []entites.Order
}

func NewMemoryOrderRepository() *MemoryOrderRepository {
	return &MemoryOrderRepository{
		mu:     sync.RWMutex{},
		orders: []entites.Order{},
	}
}

func (repo *MemoryOrderRepository) SaveOrder(order entites.Order) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.orders = append(repo.orders, order)
	return nil
}
