package service

import (
	"Go-PetStoreApp/helper"
	"Go-PetStoreApp/model/domain"
	"Go-PetStoreApp/model/web"
	"Go-PetStoreApp/repository"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	UserRepository repository.UserRepository
	DB             *sql.DB
	Validate       *validator.Validate
}

func NewUserService(userRepository repository.UserRepository, DB *sql.DB, validate *validator.Validate) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
		DB:             DB,
		Validate:       validate,
	}
}

// ChangePassword implements UserService.
func (s *UserServiceImpl) ChangePassword(ctx context.Context, request web.UserChangePasswordRequest) error {
	err := s.Validate.Struct(request)
    if err != nil {
        return err
    }

    tx, err := s.DB.Begin()
    if err != nil {
        return err
    }
    defer helper.CommitOrRollback(tx)

    user, err := s.UserRepository.FindById(ctx, tx, request.Id)
    if err != nil {
        return err
    }

    // Verify old password
    err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.OldPassword))
    if err != nil {
        return errors.New("old password is incorrect")
    }

    // Hash new password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    user.PasswordHash = string(hashedPassword)
    user.UpdatedAt = time.Now()

    SQL := "UPDATE users SET password_hash = $1, updated_at = $2 WHERE id = $3"
    _, err = tx.ExecContext(ctx, SQL, user.PasswordHash, user.UpdatedAt, user.ID)
    helper.PanicIfError(err)

    return nil
}

// Delete implements UserService.
func (s *UserServiceImpl) Delete(ctx context.Context, userId int) error {
	tx, err := s.DB.Begin()
    if err != nil {
        return err
    }
    defer helper.CommitOrRollback(tx)

    _, err = s.UserRepository.FindById(ctx, tx, userId)
    if err != nil {
        return err
    }

    s.UserRepository.Delete(ctx, tx, userId)
    return nil
}

func (s *UserServiceImpl) FindById(ctx context.Context, userId int) (web.UserResponse, error) {
    tx, err := s.DB.Begin()
    if err != nil {
        return web.UserResponse{}, err
    }
    defer helper.CommitOrRollback(tx)

    user, err := s.UserRepository.FindById(ctx, tx, userId)
    if err != nil {
        return web.UserResponse{}, err
    }

    return helper.ToUserResponse(user), nil
}

// FindAll implements UserService.
func (s *UserServiceImpl) FindAll(ctx context.Context) []web.UserResponse {
	tx, err := s.DB.Begin()
    helper.PanicIfError(err)
    defer helper.CommitOrRollback(tx)

    users := s.UserRepository.FindAll(ctx, tx)

    var userResponses []web.UserResponse
    for _, user := range users {
        userResponses = append(userResponses, helper.ToUserResponse(user))
    }

    return userResponses
}

// Login implements UserService.
func (s *UserServiceImpl) Login(ctx context.Context, request web.UserLoginRequest) (web.AuthResponse, error) {
	err := s.Validate.Struct(request)
    if err != nil {
        return web.AuthResponse{}, err
    }

    tx, err := s.DB.Begin()
    if err != nil {
        return web.AuthResponse{}, err
    }
    defer helper.CommitOrRollback(tx)

    user, err := s.UserRepository.FindByUsername(ctx, tx, request.Username)
    if err != nil {
        return web.AuthResponse{}, errors.New("invalid username or password")
    }

    // Compare password
    err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password))
    if err != nil {
        return web.AuthResponse{}, errors.New("invalid username or password")
    }

    // Generate JWT token
    token, err := helper.GenerateToken(user.ID, user.Email)
    if err != nil {
        return web.AuthResponse{}, err
    }

    return web.AuthResponse{
        Token: token,
        User:  helper.ToUserResponse(user),
    }, nil
}

// Register implements UserService.
func (s *UserServiceImpl) Register(ctx context.Context, request web.UserRegisterRequest) (web.AuthResponse, error) {
	err := s.Validate.Struct(request)
    if err != nil {
        return web.AuthResponse{}, err
    }

    tx, err := s.DB.Begin()
    if err != nil {
        return web.AuthResponse{}, err
    }
    defer helper.CommitOrRollback(tx)

    // Check if email already exists
    existingUser, _ := s.UserRepository.FindByEmail(ctx, tx, request.Email)
    if existingUser.ID != 0 {
        return web.AuthResponse{}, errors.New("email already registered")
    }

    // Check if username already exists
    existingUser, _ = s.UserRepository.FindByUsername(ctx, tx, request.Username)
    if existingUser.ID != 0 {
        return web.AuthResponse{}, errors.New("username already taken")
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
    if err != nil {
        return web.AuthResponse{}, err
    }

    now := time.Now()
    user := domain.User{
        Username:     request.Username,
        Email:        request.Email,
        PasswordHash: string(hashedPassword),
        CreatedAt:    now,
        UpdatedAt:    now,
    }

    user = s.UserRepository.Create(ctx, tx, user)

    // Generate JWT token
    token, err := helper.GenerateToken(user.ID, user.Email)
    if err != nil {
        return web.AuthResponse{}, err
    }

    return web.AuthResponse{
        Token: token,
        User:  helper.ToUserResponse(user),
    }, nil
}


// Update implements UserService.
func (s *UserServiceImpl) Update(ctx context.Context, request web.UserUpdateRequest) (web.UserResponse, error) {
	err := s.Validate.Struct(request)
    if err != nil {
        return web.UserResponse{}, err
    }

    tx, err := s.DB.Begin()
    if err != nil {
        return web.UserResponse{}, err
    }
    defer helper.CommitOrRollback(tx)

    user, err := s.UserRepository.FindById(ctx, tx, request.Id)
    if err != nil {
        return web.UserResponse{}, err
    }

    // Check if email is taken by another user
    if user.Email != request.Email {
        existingUser, _ := s.UserRepository.FindByEmail(ctx, tx, request.Email)
        if existingUser.ID != 0 && existingUser.ID != request.Id {
            return web.UserResponse{}, errors.New("email already taken")
        }
    }

    // Check if username is taken by another user
    if user.Username != request.Username {
        existingUser, _ := s.UserRepository.FindByUsername(ctx, tx, request.Username)
        if existingUser.ID != 0 && existingUser.ID != request.Id {
            return web.UserResponse{}, errors.New("username already taken")
        }
    }

    user.Username = request.Username
    user.Email = request.Email
    user.UpdatedAt = time.Now()

    user = s.UserRepository.Update(ctx, tx, user)

    return helper.ToUserResponse(user), nil
}
