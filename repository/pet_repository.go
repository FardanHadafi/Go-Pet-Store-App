package repository

import (
	"Go-PetStoreApp/model/domain"
	"context"
	"database/sql"
)

type PetRepository interface {
	Create(ctx context.Context, tx *sql.Tx, pet domain.Pet) domain.Pet
	FindById(ctx context.Context, tx *sql.Tx, id int) (domain.Pet, error)
	FindAllWithFilterByUser(ctx context.Context, tx *sql.Tx, userID, limit, offset int, species string) ([]domain.Pet, int)
	Update(ctx context.Context, tx *sql.Tx, pet domain.Pet) domain.Pet
	Delete(ctx context.Context, tx *sql.Tx, id int)
}