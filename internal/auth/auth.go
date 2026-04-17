package auth

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const claimsKey contextKey = "claims"
const guildIDKey contextKey = "guildID"

func JWTMiddleware(secret *ecdsa.PublicKey) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing authorization header", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}
				return secret, nil
			}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
			if err != nil || !token.Valid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			guildID, ok := claims["guildID"].(string)
			if !ok {
				http.Error(w, "Missing guild ID claim", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), guildIDKey, guildID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GuildIDFromContext(ctx context.Context) (string, bool) {
	guildID, ok := ctx.Value(guildIDKey).(string)
	return guildID, ok
}
