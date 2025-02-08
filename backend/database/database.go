package database

import(
	"context"
	"fmt"
	
	"github.com/jackc/pgx/v5/pgxpool"
)

func EstablishConnection() *pgxpool.Pool {
	//addr := "postgres://postgres:password@localhost:5432/pings"
	addrDocker := "postgres://postgres:password@postgres:5432/pings"
	dbpool, err := pgxpool.New(context.Background(), addrDocker)
	if err != nil {
		fmt.Printf(err.Error())
	}

	return dbpool
}
