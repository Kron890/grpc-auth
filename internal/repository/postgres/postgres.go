package postgres

import (
	"auth-grpc/internal/domain/filters"
	"auth-grpc/internal/infrastructure/storage"
	"context"
)

type RepositoryPostgres struct {
	DB *storage.DataBase
}

// TODO:...
func New(db *storage.DataBase) *RepositoryPostgres {
	return &RepositoryPostgres{DB: db}
}

// Create добавляет в базу данных логин и хеш-пароля
func (r *RepositoryPostgres) Create(ctx context.Context, user filters.UserDB) (int64, error) {
	const query = `
		INSERT INTO user_list (login, pass_hash)
		VALUES ($1, $2)
		RETURNING id
	`
	var id int64
	err := r.DB.QueryRowContext(ctx, query, user.Login, user.PassHash).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GetUser TODO:...
func (r *RepositoryPostgres) GetUser(ctx context.Context, login string) (filters.UserDB, error) {
	const query = `SELECT id, pass_hash FROM user_list WHERE login = $1`

	var u filters.UserDB
	err := r.DB.QueryRowContext(ctx, query, login).Scan(&u.ID, &u.PassHash)
	if err != nil {
		return filters.UserDB{}, err
	}

	return u, nil
}

// TODO:...
func (r *RepositoryPostgres) App(ctx context.Context, appID int) {

}

// TODO: ПРОВЕРИТЬ
func (r *RepositoryPostgres) CheckUser(ctx context.Context, login string) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM user_list WHERE login = $1)`
	var exists bool
	err := r.DB.QueryRowContext(ctx, query, login).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
