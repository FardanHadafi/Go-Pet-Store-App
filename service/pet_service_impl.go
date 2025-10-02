package service

import (
	"Go-PetStoreApp/errorsx"
	"Go-PetStoreApp/helper"
	"Go-PetStoreApp/model/domain"
	"Go-PetStoreApp/model/web"
	"Go-PetStoreApp/repository"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type PetServiceImpl struct {
	PetRepository repository.PetRepository
	DB            *sql.DB
}

func NewPetService(repo repository.PetRepository, db *sql.DB) PetService {
	return &PetServiceImpl{PetRepository: repo, DB: db}
}

func (s *PetServiceImpl) Create(ctx context.Context, req web.PetCreateRequest, userID int) (web.PetResponse, error) {
	pet := domain.Pet{
		Name:      req.Name,
		Species:   req.Species,
		Price:     req.Price,
		CreatedBy: userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return web.PetResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	created := s.PetRepository.Create(ctx, tx, pet)
	return helper.ToPetResponse(created), nil
}

func (s *PetServiceImpl) FindAllByUser(ctx context.Context, userID, page, limit int, species string) ([]web.PetResponse, int, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, 0, err
	}
	defer helper.CommitOrRollback(tx)

	offset := (page - 1) * limit
	pets, total := s.PetRepository.FindAllWithFilterByUser(ctx, tx, userID, limit, offset, species)

	res := make([]web.PetResponse, 0, len(pets))
	for _, p := range pets {
		res = append(res, helper.ToPetResponse(p))
	}
	return res, total, nil
}

func (s *PetServiceImpl) FindById(ctx context.Context, petID int, userID int) (web.PetResponse, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return web.PetResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	pet, err := s.PetRepository.FindById(ctx, tx, petID)
	if err != nil {
		return web.PetResponse{}, fmt.Errorf("%w: pet not found", errorsx.ErrNotFound)
	}
	if pet.CreatedBy != userID {
		return web.PetResponse{}, fmt.Errorf("%w: not owner", errorsx.ErrUnauthorized)
	}
	return helper.ToPetResponse(pet), nil
}

func (s *PetServiceImpl) Update(ctx context.Context, req web.PetUpdateRequest, userID int) (web.PetResponse, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return web.PetResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	pet, err := s.PetRepository.FindById(ctx, tx, req.Id)
	if err != nil {
		return web.PetResponse{}, fmt.Errorf("%w: pet not found", errorsx.ErrNotFound)
	}
	if pet.CreatedBy != userID {
		return web.PetResponse{}, fmt.Errorf("%w: not owner", errorsx.ErrUnauthorized)
	}

	pet.Name = req.Name
	pet.Species = req.Species
	pet.Price = req.Price
	pet.UpdatedAt = time.Now()

	updated := s.PetRepository.Update(ctx, tx, pet)
	return helper.ToPetResponse(updated), nil
}

func (s *PetServiceImpl) Delete(ctx context.Context, petID int, userID int) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	pet, err := s.PetRepository.FindById(ctx, tx, petID)
	if err != nil {
		return fmt.Errorf("%w: pet not found", errorsx.ErrNotFound)
	}
	if pet.CreatedBy != userID {
		return fmt.Errorf("%w: not owner", errorsx.ErrUnauthorized)
	}

	s.PetRepository.Delete(ctx, tx, petID)
	return nil
}
