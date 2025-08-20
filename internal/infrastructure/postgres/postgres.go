package postgres

import (
	"auth-grpc/internal/config"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DataBase struct {
	DB *sql.DB
}

func New(cfg *config.Config) (*DataBase, error) {

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User,
		cfg.DBPassword,
		cfg.Host,
		cfg.DBPort,
		cfg.DB,
	)

	connect, err := sql.Open("postgres", connStr)
	if err != nil {
		return &DataBase{}, err
	}

	err = connect.Ping()
	if err != nil {
		connect.Close()
		return &DataBase{}, err
	}

	table := `
CREATE TABLE IF NOT EXISTS user_list (
    id SERIAL PRIMARY KEY,
    login TEXT NOT NULL UNIQUE,
    pass_hash BYTEA NOT NULL
);`

	_, err = connect.Exec(table)
	if err != nil {
		return &DataBase{}, err
	}
	return &DataBase{DB: connect}, err
}
