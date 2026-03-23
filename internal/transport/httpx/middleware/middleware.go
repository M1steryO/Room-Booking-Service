package middleware

import (
	"context"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/domain"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/platform/security"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/transport/httpx/models"
	"net/http"
	"strings"
)

type actorContextKey struct{}

type Actor struct {
	UserID string
	Role   domain.Role
}

func AuthMiddleware(tokens *security.JWTManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorization := r.Header.Get("Authorization")
			const prefix = "Bearer "

			if !strings.HasPrefix(authorization, prefix) {
				models.WriteError(w, domain.Unauthorized("missing bearer token"))
				return
			}

			rawToken := strings.TrimPrefix(authorization, prefix)
			claims, err := tokens.Parse(rawToken)
			if err != nil {
				models.WriteError(w, domain.Unauthorized("invalid token"))
				return
			}

			ctx := context.WithValue(r.Context(), actorContextKey{}, Actor{
				UserID: claims.UserID,
				Role:   claims.Role,
			})

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ActorFromContext(ctx context.Context) Actor {
	actor, _ := ctx.Value(actorContextKey{}).(Actor)
	return actor
}
