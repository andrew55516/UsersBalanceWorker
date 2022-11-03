package db

import "github.com/jackc/pgx/v4/pgxpool"

type ServicesInstance struct {
	Db *pgxpool.Pool
}
