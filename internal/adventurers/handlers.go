package adventurers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Barms1218/nagl/internal/auth"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type AdventurerLister interface {
	ListRecruitableAdventurers(ctx context.Context, filters SearchFilters) ([]ListAdventurersResponse, error)
}

func ListRecruitableAdventurers(s AdventurerLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var filters SearchFilters
		if err := json.NewDecoder(r.Body).Decode(&filters); err != nil {
			http.Error(w, "Invalid Request Body", http.StatusBadRequest)
			return
		}

		members, err := s.ListRecruitableAdventurers(ctx, filters)
		if err != nil {
			http.Error(w, "Request Failed: %w", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(members); err != nil {
			http.Error(w, "JSON Encoding error: %w", http.StatusInternalServerError)
			return
		}
	}
}

type GuildLister interface {
	ListGuildMembers(ctx context.Context, guildID uuid.UUID, filters GuildMemberFilters) ([]ListMembersResponse, error)
}

func ListGuildMembers(g GuildLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		guildID, ok := auth.GuildIDFromContext(ctx)
		if !ok {
			http.Error(w, "Missing Guild ID Claim", http.StatusUnauthorized)
			return
		}

		var filters GuildMemberFilters
		if err := json.NewDecoder(r.Body).Decode(&filters); err != nil {
			http.Error(w, fmt.Sprintf("Bad Request Body: %v", err), http.StatusInternalServerError)
			return
		}

		adventurers, err := g.ListGuildMembers(ctx, guildID, filters)
		if err != nil {
			http.Error(w, "Adventurer list failed: %w", http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(adventurers); err != nil {
			http.Error(w, "JSON Encoding Error", http.StatusInternalServerError)
			return
		}
	}
}

type DetailsGetter interface {
	GetAdventurerDetails(ctx context.Context, id uuid.UUID) (DetailsResponse, error)
}

func GetDetails(s DetailsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		_, ok := auth.GuildIDFromContext(ctx)
		if !ok {
			http.Error(w, "Missing Guild ID Claim", http.StatusUnauthorized)
			return
		}

		idStr := chi.URLParam(r, "id")
		adventurerID, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, "Invalid adventurer ID", http.StatusBadRequest)
			return
		}

		details, err := s.GetAdventurerDetails(ctx, adventurerID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Details request failed: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(details); err != nil {
			http.Error(w, "JSON Encoding error", http.StatusInternalServerError)
			return
		}
	}
}

type Recruiter interface {
	HireAdventurer(ctx context.Context, r SetAdventurerHiredRequest) error
}

func HireAdventurer(s Recruiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		guildID, ok := auth.GuildIDFromContext(ctx)
		if !ok {
			http.Error(w, "Missing Guild ID Claim", http.StatusUnauthorized)
			return
		}

		idStr := chi.URLParam(r, "id")
		adventurerID, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, "Invalid adventurer ID", http.StatusBadRequest)
			return
		}

		request := SetAdventurerHiredRequest{
			GuildID:      guildID,
			AdventurerID: adventurerID,
		}

		if err := s.HireAdventurer(ctx, request); err != nil {
			http.Error(w, fmt.Sprintf("Request Failed: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

type UpkeepDisplayer interface {
	GetUpkeepCost(ctx context.Context, adventurerID uuid.UUID) (int32, error)
}

func GetUpkeepCost(u UpkeepDisplayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		_, ok := auth.GuildIDFromContext(ctx)
		if !ok {
			http.Error(w, "Missing Guild ID Claim", http.StatusUnauthorized)
			return
		}

		idStr := chi.URLParam(r, "id")
		adventurerID, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, "Invalid adventurer ID", http.StatusBadRequest)
			return
		}

		cost, err := u.GetUpkeepCost(ctx, adventurerID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Request Failed: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(cost); err != nil {
			http.Error(w, "JSON Encoding error", http.StatusInternalServerError)
			return
		}
	}
}
