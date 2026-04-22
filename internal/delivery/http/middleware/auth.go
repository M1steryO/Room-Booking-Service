package middleware

import (
	"context"
	"github.com/M1steryO/Room-Booking-Service/pkg/security"
	"net/http"
	"strings"

	"github.com/M1steryO/Room-Booking-Service/internal/delivery/http/models"
	"github.com/M1steryO/Room-Booking-Service/internal/domain"
)

type actorContextKey struct{}

type Actor struct {
	UserID string
	Role   domain.Role
}

func ActorFromContext(ctx context.Context) Actor {
	actor, _ := ctx.Value(actorContextKey{}).(Actor)
	return actor
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
