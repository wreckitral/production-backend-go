package post

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Route(r chi.Router, h *Handler, jwt func(http.Handler) http.Handler) {
	r.Get("/api/posts", h.List)
	r.Get("/api/posts/{id}", h.Get)

	r.With(jwt).Post("/api/posts", h.Create)
	r.With(jwt).Put("/api/posts/{id}", h.Update)
	r.With(jwt).Delete("/api/posts/{id}", h.Delete)
}
