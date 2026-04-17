package adventurers

import (
	"context"
	"encoding/json"
	"github.com/Barms1218/nagl/internal/auth"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type AdventurerStore interface {
	ListAdventurers(ctx context.Context, request GetMembersRequest) ([]GetMembersResponse, error)
	GetAdventurerDetails(ctx context.Context, id uuid.UUID) (DetailsResponse, error)
}

func ListAdventurers(s AdventurerStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var request GetMembersRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid Request Body", http.StatusBadRequest)
			return
		}

		members, err := s.ListAdventurers(ctx, request)
		if err != nil {
			http.Error(w, "Request Failed: %w", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(members); err != nil {
			http.Error(w, "JSON Encoding error: %w", http.StatusInternalServerError)
			return
		}
	}
}

func GetDetails(s AdventurerStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		id, ok := auth.GuildIDFromContext(ctx)
		if !ok {
			http.Error(w, "Missing Guild ID Claim", http.StatusUnauthorized)
			return
		}

		details, err := s.GetAdventurerDetails(ctx, id)
		if err != nil {
			http.Error(w, "Details request failed: %w", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(details); err != nil {
			http.Error(w, "JSON Encoding error: %w", http.StatusInternalServerError)
			return
		}
	}
}
