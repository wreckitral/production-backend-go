package post

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestServiceCreate_ValidatesTitle(t *testing.T) {
	svc := NewService(nil)

	_, err := svc.Create(context.Background(), uuid.New(), CreatePostRequest{
		Title: "",
		Body:  "body",
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
}
