package service

import (
	"Go-PetStoreApp/model/web"
	"context"
)

type PetService interface {
	Create(ctx context.Context, req web.PetCreateRequest, userID int) (web.PetResponse, error)
	FindAllByUser(ctx context.Context, userID, page, limit int, species string) ([]web.PetResponse, int, error)
	FindById(ctx context.Context, petID int, userID int) (web.PetResponse, error)
	Update(ctx context.Context, req web.PetUpdateRequest, userID int) (web.PetResponse, error)
	Delete(ctx context.Context, petID int, userID int) error
}