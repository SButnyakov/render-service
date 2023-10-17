package repos

import (
	"backend-api/internal/storage"
	"backend-api/internal/storage/postgres"
	"fmt"
	"github.com/lib/pq"
)

type PaymentRepository struct {
	pg *postgres.PGStorage
}

func NewPaymentRepository(pg *postgres.PGStorage) *PaymentRepository {
	return &PaymentRepository{pg: pg}
}

func (p *PaymentRepository) Create(payment storage.Payment) error {
	const fn = "postgres.repos.payments.Create"

	stmt, err := p.pg.Db.Prepare("INSERT INTO payments (date, type_id, user_id) VALUES ($1, $2, $3)")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", fn, err)
	}

	_, err = stmt.Exec(payment.DateTime, payment.TypeId, payment.UserID)
	if err != nil {
		if postgresErr, ok := err.(*pq.Error); ok && postgresErr.Code == postgres.UniqueViolationErrorCode {
			return fmt.Errorf("%s: %w", fn, storage.ErrUserExists)
		}

		return fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	return nil
}
