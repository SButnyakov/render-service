package repos

import (
	"backend-api/internal/storage"
	"backend-api/internal/storage/postgres"
	"context"
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

func (s *SubscriptionRepository) Create(subscription storage.Subscription, payment storage.Payment) error {
	const fn = "postgres.repos.subscription.Create"

	ctx := context.Background()
	tx, err := s.pg.Db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: prepare transaction: %w", fn, err)
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO payments (date, type_id, user_id) VALUES ($1, $2, $3)",
		payment.DateTime, payment.TypeId, payment.UserID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	_, err = tx.ExecContext(ctx,
		"INSERT INTO subscriptions (user_id, type_id, sub_expire_date) VALUES ($1, $2, $3)",
		subscription.UserId, subscription.TypeId, subscription.ExpireDate)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s: commit transaction: %w", fn, err)
	}

	return nil
}

func (s *SubscriptionRepository) Update(subscription storage.Subscription, payment storage.Payment) error {
	const fn = "postgres.repos.subscription.Update"

	ctx := context.Background()
	tx, err := s.pg.Db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: prepare transaction: %w", fn, err)
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO payments (date, type_id, user_id) VALUES ($1, $2, $3)",
		payment.DateTime, payment.TypeId, payment.UserID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	_, err = tx.ExecContext(ctx,
		"UPDATE subscriptions SET type_id = $2, sub_expire_date = $3 WHERE user_id = $1",
		subscription.UserId, subscription.TypeId, subscription.ExpireDate)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s: commit transaction: %w", fn, err)
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
