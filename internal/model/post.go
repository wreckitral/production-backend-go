package model

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID          uuid.UUID
	AuthorID    uuid.UUID
	AuthorEmail string
	Title       string
	Body        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
