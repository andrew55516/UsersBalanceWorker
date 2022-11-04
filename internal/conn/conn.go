package conn

import (
	"UsersBalanceWorker/pkg/helpers/e"
	"UsersBalanceWorker/pkg/helpers/pg"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

func Connection(cfg *pg.Config) (*pgxpool.Pool, error) {
	poolConfig, err := pg.NewPoolConfig(cfg)
	if err != nil {
		return nil, e.Wrap(fmt.Sprintf("%s pool config error", cfg.DbName), err)
	}

	poolConfig.MaxConns = 5

	conn, err := pg.NewConnection(poolConfig)
	if err != nil {
		return nil, e.Wrap(fmt.Sprintf("%s connection to database failed", cfg.DbName), err)
	}

	_, err = conn.Exec(context.Background(), ";")
	if err != nil {
		return nil, e.Wrap(fmt.Sprintf("%s ping failed", cfg.DbName), err)
	}
	log.Printf("%s Ping OK!\n", cfg.DbName)

	return conn, nil
}

//func Ping(db *entities.Db) error {
//	_, err := db.Users.Db.Exec(context.Background(), ";")
//	if err != nil {
//		return err
//	}
//
//	_, err = db.Services.Db.Exec(context.Background(), ";")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	_, err = db.Record.Db.Exec(context.Background(), ";")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	return nil
//}
