package repos

import (
	"backend-api/internal/storage/postgres"
	"fmt"
)

type SubscriptionTypeRepository struct {
	pg *postgres.PGStorage
}

func NewSubscriptionTypeRepository(pg *postgres.PGStorage) *SubscriptionTypeRepository {
	return &SubscriptionTypeRepository{pg: pg}
}

func (p *SubscriptionTypeRepository) GetID(name string) error {
	const fn = "postgres.repos.payment_types.GetID"

	stmt, err := p.pg.Db.Prepare("SELECT id FROM subscription_types WHERE name = $1")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", fn, err)
	}

	_, err = stmt.Exec(name)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	return nil
}
