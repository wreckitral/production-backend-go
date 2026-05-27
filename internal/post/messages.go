package post

import (
	"time"

	"github.com/google/uuid"
	"github.com/wreckitral/production-backend-go/internal/model"
)

type CreatePostRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type UpdatePostRequest struct {
	Title *string `json:"title,omitempty"`
	Body  *string `json:"body,omitempty"`
}

type PostResponse struct {
	ID          uuid.UUID `json:"id"`
	AuthorID    uuid.UUID `json:"author_id"`
	AuthorEmail string    `json:"author_email"`
	Title       string    `json:"title"`
	Body        string    `json:"body"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func toResponse(p model.Post) PostResponse {
	return PostResponse{
		ID:          p.ID,
		AuthorID:    p.AuthorID,
		AuthorEmail: p.AuthorEmail,
		Title:       p.Title,
		Body:        p.Body,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
