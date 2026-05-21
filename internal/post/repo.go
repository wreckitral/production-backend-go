package post

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
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

func (r *Repo) Create(ctx context.Context, p model.Post) (model.Post, error) {
	return model.Post{}, nil
}

func (r *Repo) GetByID(ctx context.Context, id uuid.UUID) (model.Post, error) {
	return model.Post{}, nil
}

func (r *Repo) List(ctx context.Context, limit, offset int) ([]model.Post, error) {
	return nil, nil
}

func (r *Repo) Update(ctx context.Context, p model.Post) (model.Post, error) {
	return model.Post{}, nil
}

func (r *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}
