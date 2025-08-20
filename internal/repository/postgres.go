package repository

import (
	"auth-grpc/internal/domain/filters"
	"auth-grpc/internal/infrastructure/postgres"
	"context"
)

type Repository struct {
	DB *postgres.DataBase
}

// TODO:...
func New(db *postgres.DataBase) *Repository {
	return &Repository{DB: db}
}

// Create добавляет в базу данных логин и хеш-пароля
func (r *Repository) Create(ctx context.Context, user filters.UserDB) (int64, error) {
	const query = `
		INSERT INTO user_list (login, pass_hash)
		VALUES ($1, $2)
		RETURNING id
	`
	var id int64
	err := r.DB.DB.QueryRowContext(ctx, query, user.Login, user.PassHash).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GetUser TODO:...
func (r *Repository) GetUser(ctx context.Context, login string) (filters.UserDB, error) {
	const query = `SELECT id, pass_hash FROM user_list WHERE login = $1`

	var u filters.UserDB
	err := r.DB.DB.QueryRowContext(ctx, query, login).Scan(&u.ID, &u.PassHash)
	if err != nil {
		return filters.UserDB{}, err
	}

	return u, nil
}

// TODO:...
func (r *Repository) App(ctx context.Context, appID int) {

}
