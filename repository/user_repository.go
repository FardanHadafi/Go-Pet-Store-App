package repository

import (
	"context"
	"database/sql"
	"Go-PetStoreApp/model/domain"
)

type UserRepository interface {
	Create(ctx context.Context, tx *sql.Tx, user domain.User) (domain.User, error)
	FindByEmail(ctx context.Context, tx *sql.Tx, email string) (domain.User, error)
	FindByUsername(ctx context.Context, tx *sql.Tx, username string) (domain.User, error)
	FindById(ctx context.Context, tx *sql.Tx, id int) (domain.User, error)
	FindAll(ctx context.Context, tx *sql.Tx) ([]domain.User, error)
	Update(ctx context.Context, tx *sql.Tx, user domain.User) (domain.User, error)
	Delete(ctx context.Context, tx *sql.Tx, id int) error
}