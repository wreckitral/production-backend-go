package post

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
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

func (r *Repo) Create(ctx context.Context, p model.Post) (model.Post, error) {
	const q = `
		INSERT INTO posts (author_id, title, body)
		VALUES ($1, $2, $3)
		RETURNING id, author_id, title, body, created_at, updated_at
	`

	var out model.Post
	err := r.pool.QueryRow(ctx, q, p.AuthorID, p.Title, p.Body).Scan(
		&out.ID,
		&out.AuthorID,
		&out.Title,
		&out.Body,
		&out.CreatedAt,
		&out.UpdatedAt,
	)
	if err != nil {
		return model.Post{}, fmt.Errorf("create post: %w", mapPGError(err))
	}

	return out, nil
}

func (r *Repo) GetByID(ctx context.Context, id uuid.UUID) (model.Post, error) {
	const q = `
		SELECT id, author_id, title, body, created_at, updated_at
		FROM posts
		WHERE id = $1
	`

	var out model.Post
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&out.ID,
		&out.AuthorID,
		&out.Title,
		&out.Body,
		&out.CreatedAt,
		&out.UpdatedAt,
	)

	if err != nil {
		return model.Post{}, fmt.Errorf("get post: %w", mapPGError(err))
	}

	return out, nil
}

func (r *Repo) List(ctx context.Context, limit, offset int) ([]model.Post, error) {
	const q = `
		SELECT id, author_id, title, body, created_at, updated_at
        FROM posts
        ORDER BY created_at DESC, id DESC
        LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(ctx, q, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list posts: %w", mapPGError(err))
	}

	defer rows.Close()

	posts := make([]model.Post, 0, limit)
	for rows.Next() {
		var p model.Post
		if err := rows.Scan(&p.ID, &p.AuthorID, &p.Title, &p.Body, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan post: %w", err)
		}
		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate posts: %w", err)
	}

	return posts, nil
}

func (r *Repo) Update(ctx context.Context, p model.Post) (model.Post, error) {
	const q = `
		UPDATE posts
		SET title = $1, body = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING id, author_id, title, body, created_at, updated_at
	`

	var out model.Post
	err := r.pool.QueryRow(ctx, q, p.Title, p.Body, p.ID).Scan(
		&out.ID,
		&out.AuthorID,
		&out.Title,
		&out.Body,
		&out.CreatedAt,
		&out.UpdatedAt,
	)
	if err != nil {
		return model.Post{}, fmt.Errorf("update post: %w", mapPGError(err))
	}

	return out, nil
}

func (r *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	const q = `DELETE FROM posts WHERE id = $1`
	tag, err := r.pool.Exec(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete post: %w", mapPGError(err))
	}

	if tag.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}

	return nil
}
