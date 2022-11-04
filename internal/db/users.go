package db

import (
	"UsersBalanceWorker/pkg/helpers/e"
	"context"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
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

func (i *UsersInstance) createUser(username string, balance float64) error {
	_, err := i.Db.Exec(context.Background(), "INSERT INTO users (username, balance) VALUES ($1, $2);", username, balance)
	if err != nil {
		return e.Wrap("can't create user", err)
	}

	return nil
}

func (i *UsersInstance) UpdateUserBalance(userID int, username string, value float64) error {
	rows, err := i.Db.Query(context.Background(), "SELECT balance FROM users WHERE id = $1;", userID)
	defer rows.Close()
	if err != nil {
		return e.Wrap("can't update user balance", err)
	}

	if !rows.Next() {
		if value > 0 {
			err = i.createUser(username, value)
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

	if balance+value >= 0 {
		_, err = i.Db.Exec(context.Background(), "UPDATE users SET balance = $1 WHERE id = $2;", balance+value, userID)
		if err != nil {
			return e.Wrap("can't update user balance", err)
		}
	} else {
		return e.Wrap("can't update user balance", NotEnoughMoney)
	}

	return nil
}

func (i *UsersInstance) Balance(id int) (float64, error) {
	rows, err := i.Db.Query(context.Background(), "SELECT balance FROM users WHERE id = $1;", id)
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
