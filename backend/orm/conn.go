package orm

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var dbpool *pgxpool.Pool

func Connect() {
	pool, err := pgxpool.New(context.Background(), "postgres://postgres:password@127.0.0.1:5432/postgres")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	dbpool = pool
}

func Close() {
	dbpool.Close()
}
