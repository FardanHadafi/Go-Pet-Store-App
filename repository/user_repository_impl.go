package repository

import (
	"Go-PetStoreApp/helper"
	"Go-PetStoreApp/model/domain"
	"context"
	"database/sql"
	"errors"
)

type UserRepositoryImpl struct{}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, user domain.User) domain.User {
	SQL := "insert into users(username, email, password_hash, created_at, updated_at) values($1, $2, $3, $4, $5) returning id"

	err := tx.QueryRowContext(ctx, SQL,
	user.Username,
	user.Email,
	user.PasswordHash,
	user.CreatedAt,
	user.UpdatedAt).Scan(&user.ID)

	helper.PanicIfError(err)
	return user
}

func (r *UserRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, user domain.User) domain.User {
    SQL := `UPDATE users 
            SET username = $1, email = $2, updated_at = $3 
            WHERE id = $4`
    
    _, err := tx.ExecContext(ctx, SQL, 
        user.Username, 
        user.Email, 
        user.UpdatedAt,
        user.ID,
    )
    
    helper.PanicIfError(err)
    return user
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, userId int) {
    SQL := "DELETE FROM users WHERE id = $1"
    _, err := tx.ExecContext(ctx, SQL, userId)
    helper.PanicIfError(err)
}

func (r *UserRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, userId int) (domain.User, error) {
    SQL := `SELECT id, username, email, password_hash, created_at, updated_at 
            FROM users WHERE id = $1`
    
    row := tx.QueryRowContext(ctx, SQL, userId)
    
    user := domain.User{}
    err := row.Scan(
        &user.ID, 
        &user.Username, 
        &user.Email, 
        &user.PasswordHash,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return user, errors.New("user not found")
        }
        helper.PanicIfError(err)
    }
    
    return user, nil
}

func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, tx *sql.Tx, email string) (domain.User, error) {
    SQL := `SELECT id, username, email, password_hash, created_at, updated_at 
            FROM users WHERE email = $1`
    
    row := tx.QueryRowContext(ctx, SQL, email)
    
    user := domain.User{}
    err := row.Scan(
        &user.ID, 
        &user.Username, 
        &user.Email, 
        &user.PasswordHash,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return user, errors.New("user not found")
        }
        helper.PanicIfError(err)
    }
    
    return user, nil
}

func (r *UserRepositoryImpl) FindByUsername(ctx context.Context, tx *sql.Tx, username string) (domain.User, error) {
    SQL := `SELECT id, username, email, password_hash, created_at, updated_at 
            FROM users WHERE username = $1`
    
    row := tx.QueryRowContext(ctx, SQL, username)
    
    user := domain.User{}
    err := row.Scan(
        &user.ID, 
        &user.Username, 
        &user.Email, 
        &user.PasswordHash,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return user, errors.New("user not found")
        }
        helper.PanicIfError(err)
    }
    
    return user, nil
}

func (r *UserRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []domain.User {
    SQL := `SELECT id, username, email, password_hash, created_at, updated_at 
            FROM users`
    
    rows, err := tx.QueryContext(ctx, SQL)
    helper.PanicIfError(err)
    defer rows.Close()

    var users []domain.User
    for rows.Next() {
        user := domain.User{}
        err := rows.Scan(
            &user.ID, 
            &user.Username, 
            &user.Email, 
            &user.PasswordHash,
            &user.CreatedAt,
            &user.UpdatedAt,
        )
        helper.PanicIfError(err)
        users = append(users, user)
    }
    
    err = rows.Err()
    helper.PanicIfError(err)
    
    return users
}