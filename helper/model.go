package helper

import (
	"Go-PetStoreApp/model/domain"
	"Go-PetStoreApp/model/web"
)

func ToPetResponse(pet domain.Pet) web.PetResponse {
	return web.PetResponse{
		Id: pet.ID,
		Name: pet.Name,
		Species: pet.Species,
		Price: pet.Price,
	}
}