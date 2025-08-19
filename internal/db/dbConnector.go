package db

import (
	"BankApp/internal/config"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type Database struct {
	Db *pgx.Conn
}

// Connect to PostgresSQL
// @param int, string, string, string, string - порт, адрес, имя пользователя, пароль, имя базы данных
func ConnectToDatabase(cfg *config.Config) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Name,
		cfg.Database.Password,
	)

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}
	return &Database{Db: conn}, nil
}

func (d *Database) Close() error {
	if d.Db != nil {
		return d.Db.Close(context.Background())
	}
	return nil
}

func (d *Database) CreateNew() *Queries {
	return New(d.Db)
}
