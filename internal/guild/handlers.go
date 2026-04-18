package guild

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"github.com/Barms1218/nagl/internal/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type GuildRegistrator interface {
	RegisterGuild(ctx context.Context, g GuildAuthRequest) (uuid.UUID, error)
}

func RegisterGuild(g GuildRegistrator, pk *ecdsa.PrivateKey) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var registerRequest GuildAuthRequest
		if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		id, err := g.RegisterGuild(ctx, registerRequest)
		if err != nil {
			http.Error(w, "Registration failed", http.StatusInternalServerError)
			return
		}
		token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
			"guildID": id,
			"exp":     time.Now().Add(time.Hour * 24).Unix(),
		})

		tokenStr, err := token.SignedString(pk)
		if err != nil {
			http.Error(w, "Token signing failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"token": tokenStr,
		})
	}
}

type GuildAuthenticator interface {
	EnterGuild(ctx context.Context, g GuildAuthRequest) (uuid.UUID, error)
}

func Login(g GuildAuthenticator, pk *ecdsa.PrivateKey) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var loginRequest GuildAuthRequest
		if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
			http.Error(w, "JSON Decode error", http.StatusInternalServerError)
			return
		}

		id, err := g.EnterGuild(ctx, loginRequest)
		if err != nil {
			http.Error(w, "Login Failed", http.StatusInternalServerError)
			return
		}
		token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
			"guildID": id,
			"exp":     time.Now().Add(time.Hour * 24).Unix(),
		})

		tokenStr, err := token.SignedString(pk)
		if err != nil {
			http.Error(w, "Token signing failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{
			"token": tokenStr,
		})
	}
}

type TreasuryUpdator interface {
	ChangeTreasuryAmount(ctx context.Context, request UpdateTreasuryRequest) error
}

func ChangeTreasuryAmount(t TreasuryUpdator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		guildID, ok := auth.GuildIDFromContext(ctx)
		if !ok {
			http.Error(w, "Missing Guild ID Claim", http.StatusUnauthorized)
			return
		}

		params := UpdateTreasuryRequest{
			GuildID: guildID,
		}
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "JSON Decoding error: %w", http.StatusInternalServerError)
			return
		}

		if err := t.ChangeTreasuryAmount(ctx, params); err != nil {
			http.Error(w, "Failed to update Treasury: %w", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
