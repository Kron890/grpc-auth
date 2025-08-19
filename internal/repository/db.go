package repository

import (
	"auth-grpc/infrastructure/postgres"
	"auth-grpc/internal/domain"
	"context"
)

type Repository struct {
	// db database.Database
}

// TODO:...
func New(db *postgres.Database) *Repository {
	return &Repository{}
}

// Create добавляет в базу данных логин и хеш-пароля
func (r *Repository) Creat(ctx context.Context, login string, passHash []byte) (domain.User, error) {
	return domain.User{}, nil
}

// GetUser вытаскивает из бд по login данные id,login,password
func (r *Repository) GetUser(ctx context.Context, login string) (domain.User, error) {
	return domain.User{}, nil
}

// TODO:...
func (r *Repository) App(ctx context.Context, appID int) {

}
