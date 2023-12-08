package repos

import (
	"backend-api/internal/storage"
	"backend-api/internal/storage/postgres"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type SubscriptionRepository struct {
	pg *postgres.PGStorage
}

func NewSubscriptionRepository(pg *postgres.PGStorage) *SubscriptionRepository {
	return &SubscriptionRepository{pg: pg}
}

func (s *SubscriptionRepository) Create(subscription storage.Subscription) error {
	const fn = "postgres.repos.subscription.Create"

	stmt, err := s.pg.Db.Prepare("INSERT INTO subscriptions (user_id, type_id, sub_expire_date) VALUES ($1, $2, $3)")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", fn, err)
	}

	_, err = stmt.Exec(subscription.UserId, subscription.TypeId, subscription.ExpireDate)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	return nil
}

func (s *SubscriptionRepository) Update(subscription storage.Subscription) error {
	const fn = "postgres.repos.subscription.Update"

	stmt, err := s.pg.Db.Prepare("UPDATE subscriptions SET type_id = $2, sub_expire_date = $3 WHERE user_id = $1")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", fn, err)
	}

	_, err = stmt.Exec(subscription.UserId, subscription.TypeId, subscription.ExpireDate)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	return nil
}

func (s *SubscriptionRepository) GetExpireDate(uid int64) (*time.Time, error) {
	const fn = "postgres.repos.subscription.GetExpireDate"

	stmt, err := s.pg.Db.Prepare("SELECT sub_expire_date FROM subscriptions WHERE user_id = $1")
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", fn, err)
	}

	var expireDate *time.Time

	_ = stmt.QueryRow(uid).Scan(&expireDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: execute statement: %w", fn, storage.ErrSubscriptionNotFound)
		}
		return nil, fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	return expireDate, nil
}
