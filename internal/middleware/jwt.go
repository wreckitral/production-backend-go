package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/wreckitral/production-backend-go/internal/auth"
	"github.com/wreckitral/production-backend-go/internal/platform/respond"
)

type contextKey string

const userIDKey contextKey = "userID"

func WithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func UserID(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(userIDKey).(uuid.UUID)
	return userID, ok
}

func NewJWT(secret string, issuer string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				respond.Error(w, r, http.StatusUnauthorized, "unauthorized")
				return
			}

			raw := strings.TrimPrefix(header, "Bearer ")
			claims := &auth.Claims{}

			token, err := jwt.ParseWithClaims(raw, claims, func(token *jwt.Token) (any, error) {
				if token.Method != jwt.SigningMethodHS256 {
					return nil, fmt.Errorf("unexpected signing method")
				}
				return []byte(secret), nil
			}, jwt.WithIssuer(issuer))
			if err != nil || !token.Valid {
				respond.Error(w, r, http.StatusUnauthorized, "unauthorized")
				return
			}

			userID, err := uuid.Parse(claims.Subject)
			if err != nil {
				respond.Error(w, r, http.StatusUnauthorized, "unauthorized")
				return
			}

			ctx := WithUserID(r.Context(), userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
