// Package postgres is used to perform operations against postgres databases
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v4/stdlib"
	"github.com/luna-duclos/instrumentedsql"
	"github.com/luna-duclos/instrumentedsql/opentracing"

	"pkg.dsb.dev/health"
	"pkg.dsb.dev/metrics"
	"pkg.dsb.dev/multierror"
	"pkg.dsb.dev/storage/database"
)

const (
	maxLifetime  = time.Minute * 30
	maxIdleConns = 20
	maxOpenConns = 200
)

// Open opens a connection to an SQL database using the provided DSN. Migrations are performed if the source
// is non-nil.
func Open(dsn string, migrations *database.MigrationSource) (*sql.DB, error) {
	dsn = os.ExpandEnv(dsn)
	drv := instrumentedsql.WrapDriver(
		stdlib.GetDefaultDriver(),
		instrumentedsql.WithTracer(opentracing.NewTracer(false)),
		instrumentedsql.WithOmitArgs(),
		instrumentedsql.WithOpsExcluded(instrumentedsql.OpSQLRowsNext),
	)

	conn, err := drv.OpenConnector(dsn)
	if err != nil {
		return nil, err
	}

	db := sql.OpenDB(conn)
	db.SetConnMaxLifetime(maxLifetime)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)

	if migrations != nil {
		err = database.MigrateUp(migrations, db)
		if err != nil {
			return nil, multierror.Append(err, db.Close())
		}
	}

	health.AddCheck("postgres", db.Ping)
	metrics.AddSQLStats(db)
	return db, db.Ping()
}

// CreateDatabaseWithUser creates a new user-password combination as the owner of a desired
// new database.
func CreateDatabaseWithUser(ctx context.Context, db *sql.DB, name, user, pass string) error {
	const (
		userQueryFmt = "CREATE USER %s WITH PASSWORD '%s'"
		dbQueryFmt   = "CREATE DATABASE %s WITH OWNER %s"
		permQueryFmt = "GRANT ALL PRIVILEGES ON DATABASE %s TO %s"
	)

	queries := []string{
		fmt.Sprintf(userQueryFmt, user, pass),
		fmt.Sprintf(dbQueryFmt, name, user),
		fmt.Sprintf(permQueryFmt, name, user),
	}

	return database.WithinTransaction(ctx, db, func(ctx context.Context, tx *sql.Tx) error {
		for _, q := range queries {
			if _, err := db.ExecContext(ctx, q); err != nil {
				return err
			}
		}

		return nil
	})
}
