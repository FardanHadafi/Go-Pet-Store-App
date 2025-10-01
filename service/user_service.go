package service

import (
	"Go-PetStoreApp/model/web"
	"context"
)

type UserService interface {
	Register(ctx context.Context, request web.UserRegisterRequest) (web.AuthResponse, error)
	Login(ctx context.Context, request web.UserLoginRequest) (web.AuthResponse, error)
	Update(ctx context.Context, request web.UserUpdateRequest) (web.UserResponse, error)
	ChangePassword(ctx context.Context, request web.UserChangePasswordRequest) error
	Delete(ctx context.Context, userId int) error
	FindById(ctx context.Context, userId int) (web.UserResponse, error)
	FindAll(ctx context.Context) []web.UserResponse
}