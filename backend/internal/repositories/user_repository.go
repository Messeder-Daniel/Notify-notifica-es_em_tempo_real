package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/messederdaniel/real-time-notifications/backend/internal/models"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (repository *UserRepository) Create(ctx context.Context, name, email, passwordHash string) (*models.User, error) {
	query := `
		INSERT INTO users (name, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING
			id::text,
			name,
			email,
			password_hash,
			created_at
	`

	var user models.User

	err := repository.db.QueryRow(ctx, query, name, email, passwordHash).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

func (repository *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT
			id::text,
			name,
			email,
			password_hash,
			created_at
		FROM users
		WHERE email = $1
	`

	var user models.User

	err := repository.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return &user, nil
}

func (repository *UserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	query := `
		SELECT
			id::text,
			name,
			email,
			password_hash,
			created_at
		FROM users
		WHERE id = $1
	`

	var user models.User

	err := repository.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}

	return &user, nil
}

func (repository *UserRepository) UpdateProfile(ctx context.Context, id, name, email string) (*models.User, error) {
	query := `
		UPDATE users
		SET name = $1,
			email = $2
		WHERE id = $3
		RETURNING
			id::text,
			name,
			email,
			password_hash,
			created_at
	`

	var user models.User

	err := repository.db.QueryRow(ctx, query, name, email, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to update user profile: %w", err)
	}

	return &user, nil
}

func (repository *UserRepository) UpdatePassword(ctx context.Context, id, passwordHash string) error {
	query := `
		UPDATE users
		SET password_hash = $1
		WHERE id = $2
	`

	result, err := repository.db.Exec(ctx, query, passwordHash, id)
	if err != nil {
		return fmt.Errorf("failed to update user password: %w", err)
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
