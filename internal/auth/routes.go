package auth

import (
	"github.com/go-chi/chi/v5"
)

func Route(r chi.Router, h *Handler) {
	r.Post("/api/auth/register", h.Register)
	r.Post("/api/auth/login", h.Login)
}
