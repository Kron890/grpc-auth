package storage

import (
	"auth-grpc/internal/config"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type DataBase struct {
	*sql.DB
	logs *logrus.Logger
}

func NewPostgres(cfg *config.Config, logs *logrus.Logger) (*DataBase, error) {

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresName,
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
	return &DataBase{connect, logs}, err
}

func (d *DataBase) Close() {
	err := d.DB.Close()
	if err != nil {
		d.logs.Error("Error closing the DATABASE connection: ", err)
	}
	d.logs.Info("stopping Redis server")
}
