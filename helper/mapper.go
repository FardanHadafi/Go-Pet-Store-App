package helper

import (
	"Go-PetStoreApp/model/domain"
	"Go-PetStoreApp/model/web"
)

func ToPetResponse(p domain.Pet) web.PetResponse {
	return web.PetResponse{
		Id:        p.ID,
		Name:      p.Name,
		Species:   p.Species,
		Price:     p.Price,
		OwnerId:   p.CreatedBy,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func ToUserResponse(u domain.User) web.UserResponse {
	return web.UserResponse{
		Id:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
