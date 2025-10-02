package repository

import (
	"Go-PetStoreApp/helper"
	"Go-PetStoreApp/model/domain"
	"context"
	"database/sql"
	"strconv"
)

type PetRepositoryImpl struct{}

func NewPetRepository() PetRepository {
	return &PetRepositoryImpl{}
}

func (r *PetRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, pet domain.Pet) domain.Pet {
	sql := `INSERT INTO pets (name, species, price, created_by, created_at, updated_at)
	        VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`
	err := tx.QueryRowContext(ctx, sql, pet.Name, pet.Species, pet.Price, pet.CreatedBy, pet.CreatedAt, pet.UpdatedAt).Scan(&pet.ID)
	helper.PanicIfError(err)
	return pet
}

func (r *PetRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (domain.Pet, error) {
	sql := `SELECT id, name, species, price, created_by, created_at, updated_at FROM pets WHERE id=$1`
	row := tx.QueryRowContext(ctx, sql, id)
	var pet domain.Pet
	err := row.Scan(&pet.ID, &pet.Name, &pet.Species, &pet.Price, &pet.CreatedBy, &pet.CreatedAt, &pet.UpdatedAt)
	if err != nil {
		return domain.Pet{}, err
	}
	return pet, nil
}

func (r *PetRepositoryImpl) FindAllWithFilterByUser(ctx context.Context, tx *sql.Tx, userID, limit, offset int, species string) ([]domain.Pet, int) {
	args := []interface{}{}
	where := ""
	argIndex := 1

	// owner filter
	if userID != 0 {
		where += " WHERE created_by = $" + strconv.Itoa(argIndex)
		args = append(args, userID)
		argIndex++
	}

	// species filter appended
	if species != "" {
		if where == "" {
			where = " WHERE species = $" + strconv.Itoa(argIndex)
		} else {
			where += " AND species = $" + strconv.Itoa(argIndex)
		}
		args = append(args, species)
		argIndex++
	}

	// count total
	countSQL := "SELECT COUNT(*) FROM pets" + where
	var total int
	_ = tx.QueryRowContext(ctx, countSQL, args...).Scan(&total)

	// fetch page
	args = append(args, limit, offset)
	limitIdx := argIndex
	offsetIdx := argIndex + 1
	dataSQL := "SELECT id, name, species, price, created_by, created_at, updated_at FROM pets" + where + " ORDER BY id DESC LIMIT $" + strconv.Itoa(limitIdx) + " OFFSET $" + strconv.Itoa(offsetIdx)
	rows, err := tx.QueryContext(ctx, dataSQL, args...)
	helper.PanicIfError(err)
	defer rows.Close()

	var pets []domain.Pet
	for rows.Next() {
		var p domain.Pet
		helper.PanicIfError(rows.Scan(&p.ID, &p.Name, &p.Species, &p.Price, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt))
		pets = append(pets, p)
	}
	return pets, total
}


func (r *PetRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, pet domain.Pet) domain.Pet {
	sql := `UPDATE pets SET name=$1, species=$2, price=$3, updated_at=$4 WHERE id=$5`
	_, err := tx.ExecContext(ctx, sql, pet.Name, pet.Species, pet.Price, pet.UpdatedAt, pet.ID)
	helper.PanicIfError(err)
	return pet
}

func (r *PetRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id int) {
	sql := `DELETE FROM pets WHERE id=$1`
	_, err := tx.ExecContext(ctx, sql, id)
	helper.PanicIfError(err)
}
