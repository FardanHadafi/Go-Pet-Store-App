package repository

import (
	"Go-PetStoreApp/model/domain"
	"context"
	"database/sql"
	"fmt"
)

type PetRepositoryImpl struct{}

func NewPetRepository() PetRepository {
	return &PetRepositoryImpl{}
}

func (r *PetRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, pet domain.Pet) domain.Pet {
	SQL := "INSERT INTO pets(name, species, price) VALUES($1, $2, $3) RETURNING id"

	var id int
	err := tx.QueryRowContext(ctx, SQL, pet.Name, pet.Species, pet.Price).Scan(&id)
	if err != nil {
		panic(fmt.Sprintf("Failed to create pet: %v", err))
	}

	pet.ID = id
	return pet
}