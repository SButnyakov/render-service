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

func (p *PaymentTypeRepository) GetID(name string) (int64, error) {
	const fn = "postgres.repos.payment_types.GetID"

	stmt, err := p.pg.Db.Prepare("SELECT id FROM payment_types WHERE name = $1")
	if err != nil {
		return 0, fmt.Errorf("%s: prepare statement: %w", fn, err)
	}

	var uid int64

	err = stmt.QueryRow(name).Scan(&uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("%s: execute statement: %w", fn, storage.ErrSubscriptionTypeDoesNotExists)
		}
		return 0, fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	return uid, nil
}
