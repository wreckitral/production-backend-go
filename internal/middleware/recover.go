package middleware

import (
	"log/slog"
	"net/http"

	"github.com/wreckitral/production-backend-go/internal/platform/respond"
)

func Recover(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if v := recover(); v != nil {
					log.Error("panic recovered",
						"panic", v,
						"path", r.URL.Path,
						"request_id", r.Header.Get("X-Request-ID"),
					)
					respond.Error(w, r, http.StatusInternalServerError, "internal error")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
