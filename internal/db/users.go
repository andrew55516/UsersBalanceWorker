package db

import (
	"UsersBalanceWorker/entities"
	"UsersBalanceWorker/pkg/helpers/e"
	"context"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type UsersInstance struct {
	Db *pgxpool.Pool
}

var NotEnoughMoney = errors.New("not enough money")

func (i *UsersInstance) Ping() error {
	_, err := i.Db.Exec(context.Background(), ";")
	if err != nil {
		return err
	}

	return nil
}

func (i *UsersInstance) UpdateUserBalance(cr entities.Credit) error {
	rows, err := i.Db.Query(context.Background(), "SELECT balance FROM users WHERE id = $1;", cr.UserID)
	defer rows.Close()
	if err != nil {
		return e.Wrap("can't update user balance", err)
	}

	if !rows.Next() {
		if cr.Value > 0 {
			err = i.createUser(cr.Username, cr.Value)
			if err != nil {
				return e.Wrap("can't update user balance", err)
			}
			return nil
		} else {
			return e.Wrap("can't update user balance", NotEnoughMoney)
		}
	}

	var balance float64
	err = rows.Scan(&balance)
	if err != nil {
		return e.Wrap("can't update user balance", err)
	}

	if balance+cr.Value >= 0 {
		_, err = i.Db.Exec(context.Background(), "UPDATE users SET balance = $1 WHERE id = $2;", balance+cr.Value, cr.UserID)
		if err != nil {
			return e.Wrap("can't update user balance", err)
		}
	} else {
		return e.Wrap("can't update user balance", NotEnoughMoney)
	}

	return nil
}

func (i *UsersInstance) Balance(b entities.Balance) (float64, error) {
	rows, err := i.Db.Query(context.Background(), "SELECT balance FROM users WHERE id = $1;", b.UserID)
	defer rows.Close()
	if err != nil {
		return 0, e.Wrap("can't get user balance", err)
	}

	if !rows.Next() {
		return 0, nil
	}

	var balance float64

	err = rows.Scan(&balance)
	if err != nil {
		return 0, e.Wrap("can't get user balance", err)
	}

	return balance, nil
}

func (i *UsersInstance) Users() (map[int]string, error) {
	rows, err := i.Db.Query(context.Background(), "SELECT id, username FROM users")
	defer rows.Close()
	if err != nil {
		return nil, e.Wrap("can't get users", err)
	}

	services := make(map[int]string)
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Println(err)
		}

		services[id] = name
	}

	return services, nil
}

func (i *UsersInstance) createUser(username string, balance float64) error {
	_, err := i.Db.Exec(context.Background(), "INSERT INTO users (username, balance) VALUES ($1, $2);", username, balance)
	if err != nil {
		return e.Wrap("can't create user", err)
	}

	return nil
}
