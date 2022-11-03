package db

import (
	"UsersBalanceWorker/pkg/helpers/e"
	"context"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type RecordInstance struct {
	Db *pgxpool.Pool
}

var DuplicateError = errors.New("duplicate key value")

func (i *RecordInstance) CreditRecord(userID int, value float64, status string) error {

	_, err := i.Db.Exec(context.Background(), "INSERT INTO credit_record (user_id, value, status) VALUES ($1, $2, $3);",
		userID, value, status)
	if err != nil {
		return e.Wrap("can't create credit record", err)
	}

	return nil
}

func (i *RecordInstance) ServiceRecord(orderID int, serviceID int, userID int, value float64, status string) error {
	if err := i.checkDuplicates(orderID); err != nil {
		return e.Wrap("can't create service record", err)
	}

	_, err := i.Db.Exec(context.Background(), "INSERT INTO service_record (id, service_id, user_id, value, status) VALUES ($1, $2, $3, $4, $5);",
		orderID, serviceID, userID, value, status)
	if err != nil {
		return e.Wrap("can't create service record", err)
	}

	return nil
}

func (i *RecordInstance) checkDuplicates(id int) error {
	rows, err := i.Db.Query(context.Background(), "SELECT * FROM service_record WHERE id = $1", id)
	defer rows.Close()

	if err != nil {
		return e.Wrap("can't check duplicates", err)
	}

	if rows.Next() {
		return DuplicateError
	}

	return nil
}

func (i *RecordInstance) UpdateServiceRecord(orderID int, status string) {
	_, err := i.Db.Exec(context.Background(), "UPDATE service_record SET status = $1 WHERE id = $2", status, orderID)
	if err != nil {
		log.Println(e.Wrap("can't update service record", err))
	}
}
