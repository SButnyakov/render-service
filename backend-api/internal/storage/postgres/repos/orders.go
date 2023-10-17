package repos

import (
	"backend-api/internal/storage"
	"backend-api/internal/storage/postgres"
)

type OrderRepository struct {
	pg *postgres.PGStorage
}

func NewOrderRepository(pg *postgres.PGStorage) *OrderRepository {
	return &OrderRepository{pg: pg}
}

func (o *OrderRepository) Create(order storage.Order) error {
	return nil
}
