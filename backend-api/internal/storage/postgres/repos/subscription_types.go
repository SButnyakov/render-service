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

func (st *SubscriptionTypeRepository) GetTypesMap() (map[string]int64, error) {
	const fn = "postgres.repos.subscription_types.GetTypesMap"

	stmt, err := st.pg.Db.Prepare("SELECT id, name FROM subscription_types")
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", fn, err)
	}

	rows, err := stmt.Query()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: execute statement: %w", fn, storage.ErrNoSubscriptionTypes)
		}
		return nil, fmt.Errorf("%s: execute statement: %w", fn, err)
	}
	defer rows.Close()

	types := make(map[string]int64)

	for rows.Next() {
		sType := storage.SubscriptionType{}
		err = rows.Scan(&sType.Id, &sType.Name)
		if err != nil {
			return nil, fmt.Errorf("%s: scanning rows: %w", fn, err)
		}
		types[sType.Name] = sType.Id
	}

	return types, nil
}
