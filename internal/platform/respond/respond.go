package respond

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	validation "github.com/jellydator/validation"
	"github.com/wreckitral/production-backend-go/internal/apperr"
)

const MaxJSONBodyBytes = 1 << 20 // 1 Mib

type ErrorResponse struct {
	Error     string            `json:"error"`
	Fields    map[string]string `json:"fields,omitempty"`
	RequestID string            `json:"request_id,omitempty"`
}

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func Error(w http.ResponseWriter, r *http.Request, status int, msg string) {
	JSON(w, status, ErrorResponse{
		Error:     msg,
		RequestID: r.Header.Get("X-Request-ID"),
	})
}

func DecodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	if ct := r.Header.Get("Content-Type"); ct != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != "application/json" {
			return fmt.Errorf("content-type must be application/json")
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, MaxJSONBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError

		switch {
		case errors.Is(err, io.EOF):
			return fmt.Errorf("body must not be empty")
		case errors.As(err, &syntaxErr):
			return fmt.Errorf("malformed json at byte %d", syntaxErr.Offset)
		case errors.As(err, &typeErr):
			return fmt.Errorf("invalid value for field %q", typeErr.Field)
		case strings.Contains(err.Error(), "request body too large"):
			return fmt.Errorf("body must be smaller than %d bytes", MaxJSONBodyBytes)
		default:
			return fmt.Errorf("invalid json")

		}
	}

	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return fmt.Errorf("body must contain only one json object")
	}
	return nil
}

func AppError(w http.ResponseWriter, r *http.Request, err error) {
	var validationErrs validation.Errors

	switch {
	case errors.As(err, &validationErrs):
		JSON(w, http.StatusBadRequest, ErrorResponse{
			Error:     "validation failed",
			Fields:    validationFields(validationErrs),
			RequestID: r.Header.Get("X-Request-ID"),
		})
	case errors.Is(err, apperr.ErrNotFound):
		Error(w, r, http.StatusNotFound, "not found")
	case errors.Is(err, apperr.ErrForbidden):
		Error(w, r, http.StatusForbidden, "forbidden")
	case errors.Is(err, apperr.ErrUnauthorized):
		Error(w, r, http.StatusUnauthorized, "unauthorized")
	case errors.Is(err, apperr.ErrConflict):
		Error(w, r, http.StatusConflict, "conflict")
	default:
		Error(w, r, http.StatusInternalServerError, "internal error")
	}
}

func validationFields(errs validation.Errors) map[string]string {
	fields := make(map[string]string, len(errs))
	for field, err := range errs {
		fields[field] = err.Error()
	}
	return fields
}
