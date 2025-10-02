package service

import (
	"Go-PetStoreApp/model/web"
	"context"
)

type PetService interface {
	Create(ctx context.Context, r web.PetCreateRequest, userID int) web.PetResponse
	Update(ctx context.Context, r web.PetUpdateRequest) web.PetResponse
	Delete(ctx context.Context, petId int)
	FindById(ctx context.Context, petId int) web.PetResponse
	FindAll(ctx context.Context, userID int) []web.PetResponse
}