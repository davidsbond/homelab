package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"pkg.dsb.dev/app"
	"pkg.dsb.dev/closers"
	"pkg.dsb.dev/flag"
	"pkg.dsb.dev/metrics"
	"pkg.dsb.dev/storage/database"
	"pkg.dsb.dev/storage/database/postgres"
)

func main() {
	a := app.New(
		app.WithRunner(run),
		app.WithFlags(
			&flag.String{
				Name:        "query",
				Usage:       "SQL query to execute",
				EnvVar:      "QUERY",
				Destination: &query,
				Required:    true,
			},
			&flag.String{
				Name:        "db-dsn",
				Usage:       "DSN for connecting to the database",
				EnvVar:      "DB_DSN",
				Destination: &dbDSN,
				Required:    true,
			},
		),
	)

	if err := a.Run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

var (
	dbDSN string
	query string
)

func run(ctx context.Context) error {
	metrics.Register(rowsAffected, lastInsertID)

	db, err := postgres.Open(dbDSN, nil)
	if err != nil {
		return err
	}
	defer closers.Close(db)
	defer closers.CloseFunc(metrics.Push)

	return database.WithinTransaction(ctx, db, func(ctx context.Context, tx *sql.Tx) error {
		result, err := tx.ExecContext(ctx, query)
		if err != nil {
			return err
		}

		affected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		rowsAffected.Set(float64(affected))

		insertID, err := result.LastInsertId()
		if err != nil {
			return err
		}
		lastInsertID.Set(float64(insertID))

		return nil
	})
}
