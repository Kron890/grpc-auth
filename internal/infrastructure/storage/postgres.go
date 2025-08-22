package storage

import (
	"auth-grpc/internal/config"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type DataBase struct {
	DB *sql.DB
}

func NewPostgres(cfg *config.Config) (*DataBase, error) {

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
	return &DataBase{DB: connect}, err
}

func (d *DataBase) Stop() error {
	if d.DB != nil {
		err := d.DB.Close()
		if err != nil {
			log.Printf("Ошибка при закрытии соединения с БД: %v", err)
			return err
		}
		log.Println("Соединение с PostgreSQL корректно закрыто")
	}
	return nil
}
