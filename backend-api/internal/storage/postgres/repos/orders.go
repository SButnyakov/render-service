package repos

import (
	"backend-api/internal/storage"
	"backend-api/internal/storage/postgres"
	"fmt"
)

type OrderRepository struct {
	pg *postgres.PGStorage
}

func NewOrderRepository(pg *postgres.PGStorage) *OrderRepository {
	return &OrderRepository{pg: pg}
}

func (o *OrderRepository) Create(order storage.Order) error {
	const fn = "postgres.repos.orders.Create"

	stmt, err := o.pg.Db.Prepare("INSERT INTO orders (filename, storingname, creation_date, user_id, status_id, is_deleted) VALUES ($1, $2, $3, $4, $5, $6)")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", fn, err)
	}

	_, err = stmt.Exec(order.FileName, order.StoringName, order.CreationDate, order.UserId, order.StatusId, false)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	return nil
}

func (o *OrderRepository) UpdateStatus(storingName string, uid, statusId int64) error {
	const fn = "postgres.repos.orders.UpdateStatus"

	stmt, err := o.pg.Db.Prepare("UPDATE orders SET status_id = $3 WHERE user_id = $1 AND storingname = $2")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", fn, err)
	}

	_, err = stmt.Exec(uid, storingName+".blend", statusId)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	return nil
}

func (o *OrderRepository) UpdateDownloadLink(uid int64, storingName, downloadLink string) error {
	const fn = "postgres.repos.orders.UpdateDownloadLink"

	stmt, err := o.pg.Db.Prepare("UPDATE orders SET download_link = $3 WHERE user_id = $1 AND storingname = $2")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", fn, err)
	}

	_, err = stmt.Exec(uid, storingName+".blend", downloadLink)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	return nil
}
