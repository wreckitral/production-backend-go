package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/wreckitral/production-backend-go/internal/model"
	"github.com/wreckitral/production-backend-go/internal/post"
)

func TestRepoCreateAndGet(t *testing.T) {
    cleanDB(t)
    repo := post.NewRepo(testPool)
    authorID := seedUser(t, "alice@example.com")

    in := model.Post{AuthorID: authorID, Title: "hello", Body: "world"}
    out, err := repo.Create(context.Background(), in)
    if err != nil { t.Fatalf("create: %v", err) }
    if out.ID == uuid.Nil { t.Fatal("id not populated") }

    got, err := repo.GetByID(context.Background(), out.ID)
    if err != nil { t.Fatalf("get: %v", err) }
    if got.Title != "hello" || got.Body != "world" {
        t.Fatalf("unexpected: %+v", got)
    }
}
