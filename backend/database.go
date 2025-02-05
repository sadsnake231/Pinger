package main

import(
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

func EstablishConnection() *pgxpool.Pool {
	dbpool, err := pgxpool.New(context.Background(), "postgres://postgres:password@localhost:5432/pings")
	if err != nil {
		fmt.Printf(err.Error())
	}

	return dbpool
}

