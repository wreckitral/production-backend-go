package post

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/go-chi/chi/v5"
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

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	// read authenticated user from context
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		respond.Error(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req CreatePostRequest
	if err := respond.DecodeJSON(w, r, &req); err != nil {
		respond.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}

	p, err := h.svc.Create(r.Context(), userID, req)
	if err != nil {
		respond.AppError(w, r, err)
		return
	}

	respond.JSON(w, http.StatusCreated, toResponse(p))
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		respond.Error(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}

	rawID := chi.URLParam(r, "id")
	postID, err := uuid.Parse(rawID)
	if err != nil {
		respond.Error(w, r, http.StatusBadRequest, "invalid id")
		return
	}

	var req UpdatePostRequest
	if err := respond.DecodeJSON(w, r, &req); err != nil {
		respond.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}

	p, err := h.svc.Update(r.Context(), userID, postID, req)
	if err != nil {
		respond.AppError(w, r, err)
		return
	}

	respond.JSON(w, http.StatusOK, toResponse(p))

}
