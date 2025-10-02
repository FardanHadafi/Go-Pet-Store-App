package service

import (
	"Go-PetStoreApp/exception"
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

func NewPetService(petRepository repository.PetRepository, DB *sql.DB, validate *validator.Validate) PetService {
	return &PetServiceImpl{
		PetRepository: petRepository,
		DB: DB,
		Validate: validate,
	}
}

func (s *PetServiceImpl) Create(ctx context.Context, request web.PetCreateRequest, userID int) web.PetResponse {
    pet := domain.Pet{
        Name:      request.Name,
        Species:   request.Species,
        Price:     request.Price,
        CreatedBy: userID, // ✅ owner from JWT
    }

    tx, err := s.DB.Begin()
    helper.PanicIfError(err)
    defer helper.CommitOrRollback(tx)

    pet = s.PetRepository.Create(ctx, tx, pet)

    return web.PetResponse{
        ID:      pet.ID,
        Name:    pet.Name,
        Species: pet.Species,
        Price:   pet.Price,
        OwnerID: pet.CreatedBy,
    }
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
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

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
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	s.PetRepository.Delete(ctx, tx, pet)
}

func (s *PetServiceImpl) FindById(ctx context.Context, petId int) web.PetResponse {
		tx, err := s.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	pet, err := s.PetRepository.FindById(ctx, tx, petId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	return helper.ToPetResponse(pet)
}

func (s *PetServiceImpl) FindAll(ctx context.Context, userID int) []web.PetResponse {
	tx, err := s.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	pets := s.PetRepository.FindAll(ctx, tx, userID)

	var responses []web.PetResponse
    for _, p := range pets {
        responses = append(responses, web.PetResponse{
            ID:      p.ID,
            Name:    p.Name,
            Species: p.Species,
            Price:   p.Price,
            OwnerID: p.CreatedBy,
        })
    }
	return responses
}