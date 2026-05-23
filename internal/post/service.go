package post

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	validation "github.com/jellydator/validation"
	"github.com/wreckitral/production-backend-go/internal/apperr"
	"github.com/wreckitral/production-backend-go/internal/model"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, authorID uuid.UUID, req CreatePostRequest) (model.Post, error)  {
	// validate business input
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Title, validation.Required, validation.Length(1, 200)),
		validation.Field(&req.Body, validation.Required, validation.Length(1, 50000)),
	); err != nil {
		return model.Post{}, fmt.Errorf("validate: %w", err)
	}

	// build the domain/database model
	p := model.Post{
		AuthorID: authorID,
		Title: req.Title,
		Body: req.Body,
	}

	return s.repo.Create(ctx, p)
}

func (s *Service) Update(ctx context.Context, callerID, postID uuid.UUID, req UpdatePostRequest) (model.Post, error) {
    if req.Title == nil && req.Body == nil {
        return model.Post{}, fmt.Errorf("validate: %w", validation.Errors{
            "body": errors.New("at least one field is required"),
        })
    }

    if err := validation.ValidateStruct(&req,
        validation.Field(&req.Title, validation.When(req.Title != nil, validation.Length(1, 200))),
        validation.Field(&req.Body,  validation.When(req.Body  != nil, validation.Length(1, 50000))),
    ); err != nil {
        return model.Post{}, fmt.Errorf("validate: %w", err)
    }

    existing, err := s.repo.GetByID(ctx, postID)
    if err != nil {
        return model.Post{}, err
    }

    if existing.AuthorID != callerID {
        return model.Post{}, apperr.ErrForbidden
    }

    if req.Title != nil {
        existing.Title = *req.Title
    }
    if req.Body != nil {
        existing.Body = *req.Body
    }

    return s.repo.Update(ctx, existing)
}

func (s *Service) List(ctx context.Context, limit, offset int) ([]model.Post, error) {
	return s.repo.List(ctx, limit, offset)
}
