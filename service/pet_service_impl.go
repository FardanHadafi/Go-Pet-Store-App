package service

import (
	"Go-PetStoreApp/helper"
	"Go-PetStoreApp/model/domain"
	"Go-PetStoreApp/model/web"
	"Go-PetStoreApp/repository"
	"context"
	"database/sql"

	"github.com/go-playground/validator"
)

type PetServiceImpl struct {
	PetRepository repository.PetRepository
	DB *sql.DB
	Validate *validator.Validate
}

func (s *PetServiceImpl) Create(ctx context.Context, r web.PetCreateRequest) web.PetResponse {
	// Validate the struct
	err := s.Validate.Struct(r)
	helper.PanicIfError(err)

	tx, err := s.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	pet := domain.Pet{
		Name: r.Name,
		Species: r.Species,
		Price: r.Price,
	}

	pet = s.PetRepository.Create(ctx, tx, pet)

	return helper.ToPetResponse(pet)
}

func (s *PetServiceImpl) Update(ctx context.Context, r web.PetUpdateRequest) web.PetResponse {
	// Validate the struct
	err := s.Validate.Struct(r)
	helper.PanicIfError(err)

	tx, err := s.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	// Find pets first dont be empty
	pet, err := s.PetRepository.FindById(ctx, tx, r.Id)
	helper.PanicIfError(err)

	pet.Name = r.Name
	pet.Species = r.Species
	pet.Price = r.Price

	pet = s.PetRepository.Update(ctx, tx, pet)

	return helper.ToPetResponse(pet)
}

func (s *PetServiceImpl) Delete(ctx context.Context, petId int) {
	tx, err := s.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	// Find pets first dont be empty
	pet, err := s.PetRepository.FindById(ctx, tx, petId)
	helper.PanicIfError(err)

	s.PetRepository.Delete(ctx, tx, pet)
}

func (s *PetServiceImpl) FindById(ctx context.Context, petId int) web.PetResponse {
		tx, err := s.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	pet, err := s.PetRepository.FindById(ctx, tx, petId)
	helper.PanicIfError(err)

	return helper.ToPetResponse(pet)
}

func (s *PetServiceImpl) FindAll(ctx context.Context) []web.PetResponse {
		tx, err := s.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	pets := s.PetRepository.FindAll(ctx, tx)

	var petResponses []web.PetResponse
	for _, pet := range pets {
		petResponses = append(petResponses, helper.ToPetResponse(pet))
	}
	return petResponses
}