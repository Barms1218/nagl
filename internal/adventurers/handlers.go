package adventurers

import (
	"context"
	"encoding/json"
	"github.com/Barms1218/nagl/internal/auth"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type AdventurerLister interface {
	ListAdventurers(ctx context.Context, request GetMembersRequest) ([]GetAdventurersResponse, error)
}

func ListAdventurers(s AdventurerLister) http.HandlerFunc {
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

type GuildLister interface {
	GetAdventurersByGuild(ctx context.Context, guildID uuid.UUID) ([]GetAdventurersResponse, error)
}

func ListByGuild(g GuildLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		guildID, ok := auth.GuildIDFromContext(ctx)
		if !ok {
			http.Error(w, "Missing Guild ID Claim", http.StatusUnauthorized)
			return
		}

		adventurers, err := g.GetAdventurersByGuild(ctx, guildID)
		if err != nil {
			http.Error(w, "Adventurer list failed: %w", http.StatusInternalServerError)
		}

		if err := json.NewEncoder(w).Encode(adventurers); err != nil {
			http.Error(w, "JSON Encoding Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

type StatusLister interface {
	GetAdventurersWithStatus(ctx context.Context, request AdventurersWithStatusRequest) ([]GetAdventurersResponse, error)
}

func ListByStatus(s StatusLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var request AdventurersWithStatusRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Bad Request Body", http.StatusBadRequest)
			return
		}

		adventurers, err := s.GetAdventurersWithStatus(ctx, request)
		if err != nil {
			http.Error(w, "Adventurer list failed: %w", http.StatusInternalServerError)
		}

		if err := json.NewEncoder(w).Encode(adventurers); err != nil {
			http.Error(w, "JSON Encoding Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

type DetailsGetter interface {
	GetAdventurerDetails(ctx context.Context, id uuid.UUID) (DetailsResponse, error)
}

func GetDetails(s DetailsGetter) http.HandlerFunc {
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
