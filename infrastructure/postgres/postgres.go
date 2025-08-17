package postgres

import "database/sql"

type Database struct {
	db sql.DB
}

func NewDBConnect(port string) (*Database, error) {
	return &Database{}, nil
}
