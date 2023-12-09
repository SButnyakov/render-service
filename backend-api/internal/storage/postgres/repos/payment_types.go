package repos

import (
	"backend-api/internal/storage"
	"backend-api/internal/storage/postgres"
	"database/sql"
	"errors"
	"fmt"
)

type PaymentTypeRepository struct {
	pg *postgres.PGStorage
}

func NewPaymentTypeRepository(pg *postgres.PGStorage) *PaymentTypeRepository {
	return &PaymentTypeRepository{pg: pg}
}

func (pt *PaymentTypeRepository) GetTypesMap() (map[string]int64, error) {
	const fn = "postgres.repos.payment_types.GetTypesMap"

	stmt, err := pt.pg.Db.Prepare("SELECT id, name FROM payment_types")
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", fn, err)
	}

	rows, err := stmt.Query()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: execute statement: %w", fn, storage.ErrNoPaymentTypes)
		}
		return nil, fmt.Errorf("%s: execute statement: %w", fn, err)
	}
	defer rows.Close()

	types := make(map[string]int64)

	for rows.Next() {
		pType := storage.OrderStatus{}
		err = rows.Scan(&pType.Id, &pType.Name)
		if err != nil {
			return nil, fmt.Errorf("%s: scanning rows: %w", fn, err)
		}
		types[pType.Name] = pType.Id
	}

	return types, nil
}
