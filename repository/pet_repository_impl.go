package repository

import (
	"Go-PetStoreApp/helper"
	"Go-PetStoreApp/model/domain"
	"context"
	"database/sql"
	"errors"
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

func (r *PetRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, pet domain.Pet) domain.Pet {
	SQL := "Update pets set name = $1 where id = $2"
	_, err := tx.ExecContext(ctx, SQL, pet.Name, pet.ID)
	helper.PanicIfError(err)
	return pet
}

func (r *PetRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, pet domain.Pet) {
	SQL := "Delete from pets where id = $1"
	_, err := tx.ExecContext(ctx, SQL, pet.ID)
	helper.PanicIfError(err)
}

func (r *PetRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, petId int) (domain.Pet, error) {
	SQL := "select id, name from pets where id = $1"
	row := tx.QueryRowContext(ctx, SQL, petId)

	pet := domain.Pet{}
	err := row.Scan(&pet.ID, &pet.Name)

	if err != nil {
		if err == sql.ErrNoRows {
			return pet, errors.New("pets is not found")
		}
		helper.PanicIfError(err)
	}
	return pet, nil
}

func (r *PetRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []domain.Pet {
	SQL := "select id, name from pets"
	rows, err := tx.QueryContext(ctx, SQL)
	helper.PanicIfError(err)
	defer rows.Close()

	var pets []domain.Pet
	for rows.Next() {
		pet := domain.Pet{}
		err := rows.Scan(&pet.ID, &pet.Name)
		helper.PanicIfError(err)
		pets = append(pets, pet)
	}
	return pets
} 