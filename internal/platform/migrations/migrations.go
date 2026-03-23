package migrations

import (
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	migratepgx "github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v4/stdlib"
)

//go:embed *.sql
var files embed.FS

func Up(pool *pgxpool.Pool) error {
	db := stdlib.OpenDB(*pool.Config().ConnConfig)
	defer func() {
		_ = db.Close()
	}()

	sourceDriver, err := iofs.New(files, ".")
	if err != nil {
		return fmt.Errorf("new migrate source: %w", err)
	}

	driver, err := migratepgx.WithInstance(db, &migratepgx.Config{})
	if err != nil {
		return fmt.Errorf("new migrate db driver: %w", err)
	}

	migrator, err := migrate.NewWithInstance("iofs", sourceDriver, "pgx", driver)
	if err != nil {
		return fmt.Errorf("new migrate instance: %w", err)
	}

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate up: %w", err)
	}

	return nil
}
