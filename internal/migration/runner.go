package migration

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"gorm.io/gorm"

	"workpulse/migrations"
)

// Run applies embedded SQL migrations using golang-migrate with the configured table name.
func Run(db *gorm.DB, migrationsTable string) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get sql db: %w", err)
	}

	sourceDriver, err := iofs.New(migrations.Files, ".")
	if err != nil {
		return fmt.Errorf("load embedded migrations: %w", err)
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{
		MigrationsTable: migrationsTable,
	})
	if err != nil {
		return fmt.Errorf("init migrate driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", driver)
	if err != nil {
		return fmt.Errorf("init migrate instance: %w", err)
	}

	if err := m.Up(); err != nil {
		var dirtyErr migrate.ErrDirty
		isDirty := errors.As(err, &dirtyErr)

		if !isDirty && err != migrate.ErrNoChange {
			return fmt.Errorf("apply migrations: %w", err)
		}

		if isDirty {
			version, dirty, verr := m.Version()
			if verr != nil {
				return fmt.Errorf("get dirty version: %w", verr)
			}
			if !dirty {
				return fmt.Errorf("apply migrations: %w", err)
			}

			// Roll back the version marker to the last completed migration so we can retry cleanly.
			target := int(version - 1)
			if target < 0 {
				target = 0
			}

			if ferr := m.Force(target); ferr != nil {
				return fmt.Errorf("force dirty migration version %d: %w", version, ferr)
			}

			if rerr := m.Up(); rerr != nil && rerr != migrate.ErrNoChange {
				return fmt.Errorf("retry migrations after forcing version %d: %w", target, rerr)
			}
		}
	}
	return nil
}
