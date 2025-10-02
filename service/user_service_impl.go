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

	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	UserRepository repository.UserRepository
	DB             *sql.DB
	Validate       *validator.Validate
	TokenExpiry    int
}

func NewUserService(userRepository repository.UserRepository, DB *sql.DB, validate *validator.Validate, tokenExpiry int) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
		DB:             DB,
		Validate:       validate,
		TokenExpiry:    tokenExpiry,
	}
}

func (s *UserServiceImpl) Register(ctx context.Context, request web.UserRegisterRequest) (web.AuthResponse, error) {
	if err := s.Validate.Struct(request); err != nil {
		return web.AuthResponse{}, fmt.Errorf("%w: %v", errorsx.ErrValidation, err)
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return web.AuthResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// check email/username
	if u, _ := s.UserRepository.FindByEmail(ctx, tx, request.Email); u.ID != 0 {
		return web.AuthResponse{}, fmt.Errorf("%w: email already registered", errorsx.ErrConflict)
	}
	if u, _ := s.UserRepository.FindByUsername(ctx, tx, request.Username); u.ID != 0 {
		return web.AuthResponse{}, fmt.Errorf("%w: username already taken", errorsx.ErrConflict)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return web.AuthResponse{}, err
	}

	now := time.Now()
	user := domain.User{
		Username:     request.Username,
		Email:        request.Email,
		PasswordHash: string(hashed),
		Role: "user",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	user = s.UserRepository.Create(ctx, tx, user)
	token, err := helper.GenerateToken(user.ID, user.Email, user.Role, s.TokenExpiry)
	if err != nil {
		return web.AuthResponse{}, err
	}
	return web.AuthResponse{Token: token, User: helper.ToUserResponse(user)}, nil
}

func (s *UserServiceImpl) Login(ctx context.Context, request web.UserLoginRequest) (web.AuthResponse, error) {
	if err := s.Validate.Struct(request); err != nil {
		return web.AuthResponse{}, fmt.Errorf("%w: %v", errorsx.ErrValidation, err)
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return web.AuthResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	user, err := s.UserRepository.FindByUsername(ctx, tx, request.Username)
	if err != nil {
		return web.AuthResponse{}, fmt.Errorf("%w: invalid username or password", errorsx.ErrUnauthorized)
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password)) != nil {
		return web.AuthResponse{}, fmt.Errorf("%w: invalid username or password", errorsx.ErrUnauthorized)
	}

	token, err := helper.GenerateToken(user.ID, user.Email, user.Role, s.TokenExpiry)
	if err != nil {
		return web.AuthResponse{}, err
	}
	return web.AuthResponse{Token: token, User: helper.ToUserResponse(user)}, nil
}

func (s *UserServiceImpl) RefreshToken(ctx context.Context, oldToken string) (web.AuthResponse, error) {
	claims, err := helper.ValidateToken(oldToken)
	if err != nil {
		return web.AuthResponse{}, err
	}
	// create new token
	newToken, err := helper.GenerateToken(claims.UserID, claims.Email, claims.Role, s.TokenExpiry)
	if err != nil {
		return web.AuthResponse{}, err
	}
	// fetch user for response
	tx, err := s.DB.Begin()
	if err != nil {
		return web.AuthResponse{}, err
	}
	defer helper.CommitOrRollback(tx)
	u, err := s.UserRepository.FindById(ctx, tx, claims.UserID)
	if err != nil {
		return web.AuthResponse{}, err
	}
	return web.AuthResponse{Token: newToken, User: helper.ToUserResponse(u)}, nil
}

func (s *UserServiceImpl) FindById(ctx context.Context, userId int) (web.UserResponse, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return web.UserResponse{}, err
	}
	defer helper.CommitOrRollback(tx)
	u, err := s.UserRepository.FindById(ctx, tx, userId)
	if err != nil {
		return web.UserResponse{}, fmt.Errorf("%w: user not found", errorsx.ErrNotFound)
	}
	return helper.ToUserResponse(u), nil
}

func (s *UserServiceImpl) FindAll(ctx context.Context) ([]web.UserResponse, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	users := s.UserRepository.FindAll(ctx, tx)
	res := make([]web.UserResponse, 0, len(users))
	for _, u := range users {
		res = append(res, helper.ToUserResponse(u))
	}
	return res, nil
}

func (s *UserServiceImpl) Update(ctx context.Context, req web.UserUpdateRequest) (web.UserResponse, error) {
	if err := s.Validate.Struct(req); err != nil {
		return web.UserResponse{}, fmt.Errorf("%w: %v", errorsx.ErrValidation, err)
	}
	tx, err := s.DB.Begin()
	if err != nil {
		return web.UserResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	user, err := s.UserRepository.FindById(ctx, tx, req.Id)
	if err != nil {
		return web.UserResponse{}, fmt.Errorf("%w: user not found", errorsx.ErrNotFound)
	}
	// email/username uniqueness
	if user.Email != req.Email {
		if e, _ := s.UserRepository.FindByEmail(ctx, tx, req.Email); e.ID != 0 && e.ID != req.Id {
			return web.UserResponse{}, fmt.Errorf("%w: email already taken", errorsx.ErrConflict)
		}
	}
	if user.Username != req.Username {
		if e, _ := s.UserRepository.FindByUsername(ctx, tx, req.Username); e.ID != 0 && e.ID != req.Id {
			return web.UserResponse{}, fmt.Errorf("%w: username already taken", errorsx.ErrConflict)
		}
	}
	user.Username = req.Username
	user.Email = req.Email
	user.UpdatedAt = time.Now()
	updated := s.UserRepository.Update(ctx, tx, user)
	return helper.ToUserResponse(updated), nil
}

func (s *UserServiceImpl) ChangePassword(ctx context.Context, req web.UserChangePasswordRequest) error {
	if err := s.Validate.Struct(req); err != nil {
		return fmt.Errorf("%w: %v", errorsx.ErrValidation, err)
	}
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	user, err := s.UserRepository.FindById(ctx, tx, req.Id)
	if err != nil {
		return fmt.Errorf("%w: user not found", errorsx.ErrNotFound)
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)) != nil {
		return fmt.Errorf("%w: old password incorrect", errorsx.ErrUnauthorized)
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashed)
	user.UpdatedAt = time.Now()
	_ = s.UserRepository.Update(ctx, tx, user)
	return nil
}

func (s *UserServiceImpl) Delete(ctx context.Context, userId int) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	if _, err := s.UserRepository.FindById(ctx, tx, userId); err != nil {
		return fmt.Errorf("%w: user not found", errorsx.ErrNotFound)
	}
	s.UserRepository.Delete(ctx, tx, userId)
	return nil
}
