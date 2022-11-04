package db

import (
	"UsersBalanceWorker/pkg/helpers/e"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type ServicesInstance struct {
	Db *pgxpool.Pool
}

func (i *ServicesInstance) Ping() error {
	_, err := i.Db.Exec(context.Background(), ";")
	if err != nil {
		return err
	}

	return nil
}

func (i *ServicesInstance) Services() (map[int]string, error) {
	rows, err := i.Db.Query(context.Background(), "SELECT * FROM services")
	defer rows.Close()
	if err != nil {
		return nil, e.Wrap("can't get services", err)
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
