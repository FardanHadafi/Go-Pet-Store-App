package repository

import (
	"Go-PetStoreApp/model/domain"
	"context"
	"database/sql"
	"errors"
)

type UserRepositoryImpl struct{}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, user domain.User) (domain.User, error) {
	query := `
		INSERT INTO users (username, email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at, role
	`
	row := tx.QueryRowContext(ctx, query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.CreatedAt,
		user.UpdatedAt,
	)

	err := row.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Role)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, tx *sql.Tx, email string) (domain.User, error) {
	query := `SELECT id, username, email, password_hash, role, created_at, updated_at FROM users WHERE email=$1`
	row := tx.QueryRowContext(ctx, query, email)
	var u domain.User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, sql.ErrNoRows
		}
		return domain.User{}, err
	}
	return u, nil
}

func (r *UserRepositoryImpl) FindByUsername(ctx context.Context, tx *sql.Tx, username string) (domain.User, error) {
	query := `SELECT id, username, email, password_hash, role, created_at, updated_at FROM users WHERE username=$1`
	row := tx.QueryRowContext(ctx, query, username)
	var u domain.User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, sql.ErrNoRows
		}
		return domain.User{}, err
	}
	return u, nil
}

func (r *UserRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (domain.User, error) {
	query := `SELECT id, username, email, password_hash, role, created_at, updated_at FROM users WHERE id=$1`
	row := tx.QueryRowContext(ctx, query, id)
	var u domain.User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, sql.ErrNoRows
		}
		return domain.User{}, err
	}
	return u, nil
}


func (r *UserRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]domain.User, error) {
	query := `SELECT id, username, email, password_hash, role, created_at, updated_at FROM users`
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, user domain.User) (domain.User, error) {
	query := `UPDATE users SET username=$1, email=$2, updated_at=$3 WHERE id=$4 RETURNING updated_at`
	err := tx.QueryRowContext(ctx, query, user.Username, user.Email, user.UpdatedAt, user.ID).Scan(&user.UpdatedAt)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	query := `DELETE FROM users WHERE id=$1`
	_, err := tx.ExecContext(ctx, query, id)
	return err
}
