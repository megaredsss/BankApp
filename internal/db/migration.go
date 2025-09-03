package db

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(conn *Database) error {
	dsn := fmt.Sprintf("pgx5://%s:%s@%s:%d/%s?sslmode=disable",
		conn.Db.Config().User,
		conn.Db.Config().Password,
		conn.Db.Config().Host,
		conn.Db.Config().Port,
		conn.Db.Config().Database,
	)
	m, err := migrate.New(
		"file://internal/db/migrations",
		dsn)
	if err != nil {
		return err
	}
	defer m.Close()
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
