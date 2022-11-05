package db

import (
	"UsersBalanceWorker/entities"
	"UsersBalanceWorker/pkg/helpers/e"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"sync"
	"time"
)

const (
	DefComFrom = "transfer from user"
	DefComTo   = "transfer to user"
)

type RecordInstance struct {
	Db *pgxpool.Pool
}

type Update struct {
	ServiceID int
	UserID    int
	Value     float64
}

type AccountUnit struct {
	ServiceID int
	Value     float64
}

type CreditUnit struct {
	Value float64
	Time  time.Time
}

type TransferUnit struct {
	Value   float64
	Time    time.Time
	Comment string
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

func (i *RecordInstance) CreditRecord(cr entities.Credit, status string) error {

	_, err := i.Db.Exec(context.Background(), "INSERT INTO credit_record (user_id, value, status) VALUES ($1, $2, $3);",
		cr.UserID, cr.Value, status)
	if err != nil {
		return e.Wrap("can't create credit record", err)
	}

	return nil
}

func (i *RecordInstance) ServiceRecord(s entities.Service, status string) error {
	if err := i.checkDuplicates(s.OrderID); err != nil {
		return e.Wrap("can't create service record", err)
	}

	_, err := i.Db.Exec(context.Background(), "INSERT INTO service_record (id, service_id, user_id, value, status) VALUES ($1, $2, $3, $4, $5);",
		s.OrderID, s.ServiceID, s.UserID, s.Cost, status)
	if err != nil {
		return e.Wrap("can't create service record", err)
	}

	return nil
}

func (i *RecordInstance) TransferRecord(t entities.Transfer, comment string, status string) error {
	_, err := i.Db.Exec(context.Background(), "INSERT INTO transfer_record (user_from_id, user_to_id, value, comment, status) VALUES ($1, $2, $3, $4, $5);",
		t.UserFromID, t.UserToID, t.Value, comment, status)
	if err != nil {
		return e.Wrap("can't create transfer record", err)
	}

	return nil
}

func (i *RecordInstance) UpdateServiceRecord(orderID int, status string) (*Update, error) {
	rows, err := i.Db.Query(context.Background(), "SELECT service_id, user_id, value, status FROM service_record WHERE id = $1", orderID)
	defer rows.Close()

	if err != nil {
		return nil, e.Wrap("can't update service record", err)
	}

	var upd Update
	var stat string

	if rows.Next() {
		if err := rows.Scan(&upd.ServiceID, &upd.UserID, &upd.Value, &stat); err != nil {
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

func (i *RecordInstance) CreditHistory(UserID int) ([]CreditUnit, error) {
	rows, err := i.Db.Query(context.Background(), "SELECT value, time FROM credit_record WHERE user_id = $1 AND status = 'ok';", UserID)
	defer rows.Close()

	if err != nil {
		return nil, e.Wrap("can't get credit history", err)
	}

	history := make([]CreditUnit, 0)

	for rows.Next() {
		var unit CreditUnit
		err = rows.Scan(&unit.Value, &unit.Time)
		if err != nil {
			log.Println(err)
		}

		history = append(history, unit)
	}

	return history, nil
}

func (i *RecordInstance) TransferHistory(UserID int) ([]TransferUnit, []TransferUnit, error) {
	rowsTo, err := i.Db.Query(context.Background(), "SELECT user_to_id, value, comment, time FROM transfer_record WHERE user_from_id = $1 AND status = 'ok';", UserID)
	defer rowsTo.Close()
	if err != nil {
		return nil, nil, err
	}

	rowsFrom, err := i.Db.Query(context.Background(), "SELECT user_from_id, value, comment, time FROM transfer_record WHERE user_to_id = $1 AND status = 'ok';", UserID)
	defer rowsFrom.Close()
	if err != nil {
		return nil, nil, err
	}

	var wg sync.WaitGroup

	transferFrom := make([]TransferUnit, 0)
	transferTo := make([]TransferUnit, 0)

	wg.Add(2)

	go func() {
		defer wg.Done()
		for rowsTo.Next() {
			var unit TransferUnit
			var id int
			err := rowsTo.Scan(&id, &unit.Value, &unit.Comment, &unit.Time)
			if err != nil {
				log.Println(err)
			}

			if unit.Comment == "" {
				unit.Comment = fmt.Sprintf("%s #%d", DefComTo, id)
			}

			transferTo = append(transferTo, unit)
		}
	}()

	go func() {
		defer wg.Done()
		for rowsFrom.Next() {
			var unit TransferUnit
			var id int
			err := rowsFrom.Scan(&id, &unit.Value, &unit.Comment, &unit.Time)
			if err != nil {
				log.Println(err)
			}

			if unit.Comment == "" {
				unit.Comment = fmt.Sprintf("%s #%d", DefComFrom, id)
			}

			transferFrom = append(transferFrom, unit)
		}
	}()

	wg.Wait()

	return transferTo, transferFrom, nil
}
