package repository

import (
	"auth-grpc/infrastructure/postgres"
	"auth-grpc/internal/domain"
	"context"
)

type Repository struct {
	// db database.Database
}

func New(db *postgres.Database) (*Repository, error) {
	return &Repository{}, nil
}

func (r *Repository) SaveUser(ctx context.Context, login string, passHash []byte) (domain.User, error) {
	return domain.User{}, nil
}

func (r *Repository) UserAuth(ctx context.Context, login string, passHash []byte) (domain.User, error) {
	return domain.User{}, nil
}

func (r *Repository) App(ctx context.Context, appID int) {

}
