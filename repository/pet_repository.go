package repository

import (
	"Go-PetStoreApp/model/domain"
	"context"
	"database/sql"
)

type PetRepository interface {
	// CreatePet - Create a new pet
	Create(ctx context.Context, tx *sql.Tx, pet domain.Pet) domain.Pet
	// UpdatePet - Update an existing pet
	Update(ctx context.Context, tx *sql.Tx, pet domain.Pet) domain.Pet
	// DeletePet - Delete a pet by ID
	Delete(ctx context.Context, tx *sql.Tx, pet domain.Pet)
	// GetPetById - Get a pet by ID
	FindById(ctx context.Context, tx *sql.Tx, petId int) (domain.Pet, error)
	// GetAllPets - Get all pets
	FindAll(ctx context.Context, tx *sql.Tx) []domain.Pet
}