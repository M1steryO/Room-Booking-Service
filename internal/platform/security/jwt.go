package security

import (
	"fmt"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	UserID string      `json:"user_id"`
	Role   domain.Role `json:"role"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTManager(secret string, ttl time.Duration) *JWTManager {
	return &JWTManager{
		secret: []byte(secret),
		ttl:    ttl,
	}
}

func (m *JWTManager) Issue(userID string, role domain.Role, now time.Time) (string, error) {
	claims := TokenClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *JWTManager) Parse(raw string) (TokenClaims, error) {
	token, err := jwt.ParseWithClaims(raw, &TokenClaims{}, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return m.secret, nil
	})
	if err != nil {
		return TokenClaims{}, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return TokenClaims{}, fmt.Errorf("invalid token")
	}

	if claims.UserID == "" || !claims.Role.Valid() {
		return TokenClaims{}, fmt.Errorf("invalid token claims")
	}

	return *claims, nil
}
