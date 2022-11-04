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

type Update struct {
	UserID int
	Value  float64
}

type AccountUnit struct {
	ServiceID int
	Value     float64
}

var DuplicateError = errors.New("duplicates: order with this id already exists")
var DoesNotExistError = errors.New("this order does not exist")
var DoneOrderError = errors.New("this order has already been done")

func (i *RecordInstance) Ping() error {
	_, err := i.Db.Exec(context.Background(), ";")
	if err != nil {
		return err
	}

	return nil
}

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

func (i *RecordInstance) TransferRecord(userFromID int, userToID int, value float64, status string) error {
	_, err := i.Db.Exec(context.Background(), "INSERT INTO transfer_record (user_from_id, user_to_id, value, status) VALUES ($1, $2, $3, $4);",
		userFromID, userToID, value, status)
	if err != nil {
		return e.Wrap("can't create transfer record", err)
	}

	return nil
}

func (i *RecordInstance) UpdateServiceRecord(orderID int, status string) (*Update, error) {
	rows, err := i.Db.Query(context.Background(), "SELECT value, user_id, status FROM service_record WHERE id = $1", orderID)
	defer rows.Close()

	if err != nil {
		return nil, e.Wrap("can't update service record", err)
	}

	var upd Update
	var stat string

	if rows.Next() {
		if err := rows.Scan(&upd.Value, &upd.UserID, &stat); err != nil {
			return nil, e.Wrap("can't update service record", err)
		}
	} else {
		return nil, e.Wrap("can't update service record", DoesNotExistError)
	}

	if stat != "in process" {
		return nil, e.Wrap("can't update service record", DoneOrderError)
	}

	_, err = i.Db.Exec(context.Background(), "UPDATE service_record SET (status, time) = ($1, 'now()') WHERE id = $2", status, orderID)
	if err != nil {
		return nil, e.Wrap("can't update service record", err)
	}

	return &upd, nil
}

func (i *RecordInstance) Accounting(from string, to string) ([]AccountUnit, error) {
	rows, err := i.Db.Query(context.Background(), "SELECT service_id, value FROM service_record WHERE (time BETWEEN $1 AND $2) AND (status = 'ok');", from, to)
	defer rows.Close()
	if err != nil {
		return nil, e.Wrap("can't get service record", err)
	}

	units := make([]AccountUnit, 0)
	for rows.Next() {
		var unit AccountUnit
		err = rows.Scan(&unit.ServiceID, &unit.Value)
		if err != nil {
			log.Println(err)
		}

		units = append(units, unit)
	}

	return units, nil
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
