package postgres

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"

	pkgdb "pkg.dsb.dev/storage/database"
)

// NewDirectoryMigration creates a migration from the provided directory path.
// The parameter must be a directory.
func NewDirectoryMigration(dir string) *pkgdb.MigrationSource {
	return pkgdb.NewDirectoryMigration(newPostgresDBDriver, dir)
}

// NewBindataMigration creates a migration from the provided bindata asset.
func NewBindataMigration(b *bindata.AssetSource) *pkgdb.MigrationSource {
	return pkgdb.NewBindataMigration(newPostgresDBDriver, b)
}

func newPostgresDBDriver(db *sql.DB) (database.Driver, error) {
	return postgres.WithInstance(db, new(postgres.Config))
}
