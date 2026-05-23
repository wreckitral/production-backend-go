package auth

import (
	"net/http"

	"github.com/wreckitral/production-backend-go/internal/platform/respond"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := respond.DecodeJSON(w, r, &req); err != nil {
		respond.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}

	u, err := h.svc.Register(r.Context(), req)
	if err != nil {
		respond.AppError(w, r, err)
		return
	}

	respond.JSON(w, http.StatusCreated, UserResponse{
		ID:    u.ID,
		Email: u.Email,
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := respond.DecodeJSON(w, r, &req); err != nil {
		respond.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}

	loginResponse, err := h.svc.Login(r.Context(), req)
	if err != nil {
		respond.AppError(w, r, err)
		return
	}

	respond.JSON(w, http.StatusOK, loginResponse)
}
