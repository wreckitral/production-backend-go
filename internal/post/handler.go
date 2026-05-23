package post

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/wreckitral/production-backend-go/internal/middleware"
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

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	limit, offer, err := parsePagination(r)
	if err != nil {
		respond.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}

	posts, err := h.svc.List(r.Context(), limit, offer)
	if err != nil {
		respond.AppError(w, r, err)
		return
	}

	out := make([]PostResponse, 0, len(posts))
	for _, p := range posts {
		out = append(out, toResponse(p))
	}

	respond.JSON(w, http.StatusOK, out)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	rawID := chi.URLParam(r, "id")
	postID, err := uuid.Parse(rawID)
	if err != nil {
		respond.Error(w, r, http.StatusBadRequest, "invalid id")
		return
	}

	post, err := h.svc.Get(r.Context(), postID)
	if err != nil {
		respond.AppError(w, r, err)
		return
	}

	respond.JSON(w, http.StatusOK, toResponse(post))
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	rawID := chi.URLParam(r, "id")
	postID, err := uuid.Parse(rawID)
	if err != nil {
		respond.Error(w, r, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.svc.Delete(r.Context(), postID); err != nil {
		respond.AppError(w, r, err)
		return
	}

	respond.JSON(w, http.StatusNoContent, nil)
}

func parsePagination(r *http.Request) (limit int, offset int, err error) {
	q := r.URL.Query()

	limit = 20
	if raw := q.Get("limit"); raw != "" {
		limit, err = strconv.Atoi(raw)
		if err != nil {
			return 0, 0, fmt.Errorf("limit must be an integer")
		}
	}
	if limit < 1 || limit > 100 {
		return 0, 0, fmt.Errorf("limit must be between 1 and 100")
	}

	offset = 0
	if raw := q.Get("offset"); raw != "" {
		offset, err = strconv.Atoi(raw)
		if err != nil {
			return 0, 0, fmt.Errorf("offset must be an integer")
		}
	}
	if offset < 0 {
		return 0, 0, fmt.Errorf("offset must be greater than or equal to 0")
	}

	return limit, offset, nil
}
