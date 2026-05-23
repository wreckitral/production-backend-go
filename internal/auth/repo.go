package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wreckitral/production-backend-go/internal/apperr"
	"github.com/wreckitral/production-backend-go/internal/model"
)

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{
		pool: pool,
	}
}

func mapPGError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return apperr.ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return apperr.ErrConflict
	}

	return err
}

func (r *Repo) CreateUser(ctx context.Context, u model.User) (model.User, error) {
	const q = `
		INSERT INTO users (email, password_hash)
		VALUES (lower($1), $2)
		RETURNING id, email, password_hash, created_at;
	`

	var out model.User
	err := r.pool.QueryRow(ctx, q, u.Email, u.PasswordHash).Scan(
		&out.ID,
		&out.Email,
		&out.PasswordHash,
		&out.CreatedAt,
	)
	if err != nil {
		return model.User{}, fmt.Errorf("create user: %w", mapPGError(err))
	}

	return out, nil
}

func (r *Repo) GetByEmail(ctx context.Context, email string) (model.User, error) {
	const q = `
		SELECT id, email, password_hash, created_at
		FROM users
		WHERE email = $1;
	`

	var out model.User
	err := r.pool.QueryRow(ctx, q, email).Scan(
		&out.ID,
		&out.Email,
		&out.PasswordHash,
		&out.CreatedAt,
	)

	if err != nil {
		return model.User{}, fmt.Errorf("get user: %w", mapPGError(err))
	}

	return out, nil
}
