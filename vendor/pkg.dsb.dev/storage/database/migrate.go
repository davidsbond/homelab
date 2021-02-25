package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/source"

	// Driver for file based migrations.
	_ "github.com/golang-migrate/migrate/v4/source/file"

	// Driver for bindata migrations.
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
)

// ErrNotADirectory is returned if attempting to migrate without a valid directory.
var ErrNotADirectory = errors.New("invalid migration directory")

// The MigrationSource type is used as a source of migration info for performing migration actions
// against a database.
type MigrationSource struct {
	name           string
	sourceDriver   func() (source.Driver, error)
	databaseDriver func(db *sql.DB) (database.Driver, error)
}

// Migrate performs the migration action given a source driver.
func Migrate(action func(*migrate.Migrate) error, source *MigrationSource, db *sql.DB) error {
	dbDriver, err := source.databaseDriver(db)
	if err != nil {
		return err
	}

	d, err := source.sourceDriver()
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance(source.name, d, "", dbDriver)
	if err != nil {
		return err
	}

	err = action(m)
	switch {
	case errors.Is(err, migrate.ErrNoChange):
		return nil
	case err != nil:
		return fmt.Errorf("could not perform migration: %w", err)
	default:
		return nil
	}
}

// MigrateUp applies all migrations from the migration source that haven't been applied yet.
func MigrateUp(s *MigrationSource, db *sql.DB) error {
	action := func(m *migrate.Migrate) error { return m.Up() }

	return Migrate(action, s, db)
}

// MigrateDown reverts all migrations from the migration source that have already been applied.
func MigrateDown(s *MigrationSource, db *sql.DB) error {
	action := func(m *migrate.Migrate) error { return m.Down() }

	return Migrate(action, s, db)
}

// NewDirectoryMigration creates a migration from the provided directory path.
// The parameter must be a directory.
func NewDirectoryMigration(dbDriver func(db *sql.DB) (database.Driver, error), dir string) *MigrationSource {
	return &MigrationSource{
		name: "migration-dir",
		sourceDriver: func() (source.Driver, error) {
			ok, err := isDirectory(dir)
			if err != nil {
				return nil, err
			}

			if !ok {
				return nil, ErrNotADirectory
			}

			return source.Open(pathToFileURL(dir))
		},
		databaseDriver: dbDriver,
	}
}

// NewBindataMigration creates a migration from the provided bindata asset.
func NewBindataMigration(dbDriver func(db *sql.DB) (database.Driver, error), b *bindata.AssetSource) *MigrationSource {
	return &MigrationSource{
		name: "bindata",
		sourceDriver: func() (source.Driver, error) {
			return bindata.WithInstance(b)
		},
		databaseDriver: dbDriver,
	}
}

func isDirectory(name string) (bool, error) {
	info, err := os.Stat(name)
	if err != nil {
		return false, err
	}

	return info.IsDir(), nil
}

func pathToFileURL(path string) string {
	if !filepath.IsAbs(path) {
		var err error
		path, err = filepath.Abs(path)
		if err != nil {
			return ""
		}
	}

	return fmt.Sprintf("file://%s", filepath.ToSlash(path))
}
