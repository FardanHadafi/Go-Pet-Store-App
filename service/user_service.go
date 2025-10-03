package service

import (
	"Go-PetStoreApp/model/web"
	"context"
)

type UserService interface {
	Register(ctx context.Context, req web.UserRegisterRequest) (web.AuthResponse, error)
	Login(ctx context.Context, req web.UserLoginRequest) (web.AuthResponse, error)
	RefreshToken(ctx context.Context, oldToken string) (web.AuthResponse, error)
	FindById(ctx context.Context, id int) (web.UserResponse, error)
	FindAll(ctx context.Context) ([]web.UserResponse, error)
	Update(ctx context.Context, id int, req web.UserUpdateRequest) (web.UserResponse, error)
	ChangePassword(ctx context.Context, req web.UserChangePasswordRequest) error
	Delete(ctx context.Context, id int) error
}