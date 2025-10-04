package service

import (
	"Go-PetStoreApp/errorsx"
	"Go-PetStoreApp/helper"
	"Go-PetStoreApp/model/domain"
	"Go-PetStoreApp/model/web"
	"Go-PetStoreApp/repository"
	"context"
	"database/sql"
	"errors"
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
	// validate request
	if err := s.Validate.Struct(request); err != nil {
		return web.AuthResponse{}, fmt.Errorf("%w: %v", errorsx.ErrValidation, err)
	}

	// start transaction
	tx, err := s.DB.Begin()
	if err != nil {
		return web.AuthResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// set role (default: user)
	role := "user"
	if request.Role != "" {
		role = request.Role
	}

	// check email/username uniqueness
	if u, _ := s.UserRepository.FindByEmail(ctx, tx, request.Email); u.ID != 0 {
		return web.AuthResponse{}, fmt.Errorf("%w: email already registered", errorsx.ErrConflict)
	}
	if u, _ := s.UserRepository.FindByUsername(ctx, tx, request.Username); u.ID != 0 {
		return web.AuthResponse{}, fmt.Errorf("%w: username already taken", errorsx.ErrConflict)
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return web.AuthResponse{}, err
	}

	// build user object
	now := time.Now()
	user := domain.User{
		Username:     request.Username,
		Email:        request.Email,
		PasswordHash: string(hashedPassword),
		Role:         role,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// save to DB
	createdUser, err := s.UserRepository.Create(ctx, tx, user)
	if err != nil {
		return web.AuthResponse{}, err
	}

	// generate JWT token
	token, err := helper.GenerateToken(createdUser.ID, createdUser.Email, createdUser.Username, createdUser.Role, s.TokenExpiry)
	if err != nil {
		return web.AuthResponse{}, err
	}

	return web.AuthResponse{
		Token: token,
		User:  helper.ToUserResponse(createdUser),
	}, nil
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

	token, err := helper.GenerateToken(user.ID, user.Email, user.Username, user.Role, s.TokenExpiry)
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
	newToken, err := helper.GenerateToken(claims.UserID, claims.Email, claims.Username, claims.Role, s.TokenExpiry)
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

func (s *UserServiceImpl) FindById(ctx context.Context, id int) (web.UserResponse, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return web.UserResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	user, err := s.UserRepository.FindById(ctx, tx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return web.UserResponse{}, fmt.Errorf("%w: user not found", errorsx.ErrNotFound)
		}
		return web.UserResponse{}, err
	}

	return helper.ToUserResponse(user), nil
}


func (s *UserServiceImpl) FindAll(ctx context.Context) ([]web.UserResponse, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	users, err := s.UserRepository.FindAll(ctx, tx)
	if err != nil {
		return nil, err
	}

	res := make([]web.UserResponse, 0, len(users))
	for _, u := range users {
		res = append(res, helper.ToUserResponse(u))
	}
	return res, nil
}


func (s *UserServiceImpl) Update(ctx context.Context, id int, request web.UserUpdateRequest) (web.UserResponse, error) {
	if err := s.Validate.Struct(request); err != nil {
		return web.UserResponse{}, fmt.Errorf("%w: %v", errorsx.ErrValidation, err)
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return web.UserResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Find user first
	user, err := s.UserRepository.FindById(ctx, tx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return web.UserResponse{}, fmt.Errorf("%w: user not found", errorsx.ErrNotFound)
		}
		return web.UserResponse{}, err
	}

	// Update fields
	user.Username = request.Username
	user.Email = request.Email
	user.UpdatedAt = time.Now()

	updatedUser, err := s.UserRepository.Update(ctx, tx, user)
	if err != nil {
		return web.UserResponse{}, err
	}

	return helper.ToUserResponse(updatedUser), nil
}

func (s *UserServiceImpl) ChangePassword(ctx context.Context, req web.UserChangePasswordRequest) error {
    tx, err := s.DB.Begin()
    if err != nil {
        return err
    }
    defer helper.CommitOrRollback(tx)

    user, err := s.UserRepository.FindById(ctx, tx, req.Id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return fmt.Errorf("%w: user not found", errorsx.ErrNotFound)
        }
        return err
    }

    // compare old password
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
        return fmt.Errorf("%w: invalid old password", errorsx.ErrUnauthorized)
    }

    // hash new password
    hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    user.PasswordHash = string(hashed)
    user.UpdatedAt = time.Now()

    _, err = s.UserRepository.Update(ctx, tx, user)
    return err
}

func (s *UserServiceImpl) Delete(ctx context.Context, id int) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	// Ensure user exists
	_, err = s.UserRepository.FindById(ctx, tx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: user not found", errorsx.ErrNotFound)
		}
		return err
	}

	// Delete using repository
	if err := s.UserRepository.Delete(ctx, tx, id); err != nil {
		return err
	}

	return nil
}
