package repos

import (
	"backend-api/internal/storage"
	"backend-api/internal/storage/postgres"
	"database/sql"
	"errors"
	"fmt"
)

type SubscriptionTypeRepository struct {
	pg *postgres.PGStorage
}

func NewSubscriptionTypeRepository(pg *postgres.PGStorage) *SubscriptionTypeRepository {
	return &SubscriptionTypeRepository{pg: pg}
}

func (p *SubscriptionTypeRepository) GetID(name string) (int64, error) {
	const fn = "postgres.repos.subscription_types.GetID"

	stmt, err := p.pg.Db.Prepare("SELECT id FROM subscription_types WHERE name = $1")
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
