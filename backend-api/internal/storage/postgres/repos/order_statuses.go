package repos

import (
	"backend-api/internal/storage"
	"backend-api/internal/storage/postgres"
	"database/sql"
	"errors"
	"fmt"
)

type OrderStatusesRepository struct {
	pg *postgres.PGStorage
}

func NewOrderStatusesRepository(pg *postgres.PGStorage) *OrderStatusesRepository {
	return &OrderStatusesRepository{pg: pg}
}

func (os *OrderStatusesRepository) GetStatusesMap() (map[string]int64, error) {
	const fn = "postgres.repos.order_statuses.GetStatusesMap"

	stmt, err := os.pg.Db.Prepare("SELECT id, name FROM order_statuses")
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", fn, err)
	}

	rows, err := stmt.Query()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: execute statement: %w", fn, storage.ErrNoOrderStatuses)
		}
		return nil, fmt.Errorf("%s: execute statement: %w", fn, err)
	}
	defer rows.Close()

	statuses := make(map[string]int64)

	for rows.Next() {
		status := storage.OrderStatus{}
		err = rows.Scan(&status.Id, &status.Name)
		if err != nil {
			return nil, fmt.Errorf("%s: scanning rows: %w", fn, err)
		}
		statuses[status.Name] = status.Id
	}

	return statuses, nil
}
