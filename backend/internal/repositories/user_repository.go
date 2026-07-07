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
