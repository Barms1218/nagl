package contracts

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Barms1218/nagl/internal/auth"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ContractClaimer interface {
	ClaimContract(ctx context.Context, c ContractClaimRequest) error
}

func ClaimContract(s ContractClaimer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		guildID, ok := auth.GuildIDFromContext(ctx)
		if !ok {
			http.Error(w, "Missing Guild ID Claim", http.StatusUnauthorized)
			return
		}

		contractID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid Contract ID: %v", err), http.StatusInternalServerError)
			return
		}

		claimRequest := ContractClaimRequest{
			ContractID: contractID,
			GuildID:    guildID,
		}

		if err := json.NewDecoder(r.Body).Decode(&claimRequest); err != nil {
			http.Error(w, fmt.Sprintf("Invalid Request Body: %v", err), http.StatusBadRequest)
			return
		}

		if err := s.ClaimContract(ctx, claimRequest); err != nil {
			http.Error(w, "", http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	}
}

type ContractStarter interface {
	StartContract(ctx context.Context, c SetContractStatusRequest) error
}

func StartContract(s ContractStarter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var request SetContractStatusRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, fmt.Sprintf("Bad Request Body: %v", err), http.StatusBadRequest)
			return
		}

		if err := s.StartContract(ctx, request); err != nil {
			http.Error(w, fmt.Sprintf("Contract Start Failed: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

type AvailableContractLister interface {
	ListAvailableContracts(ctx context.Context, filter SearchFilters) ([]ListContractsResponse, error)
}

func ListAvailableContracts(s AvailableContractLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var filters SearchFilters
		if err := json.NewDecoder(r.Body).Decode(&filters); err != nil {
			http.Error(w, fmt.Sprintf("Bad Request: %v", err), http.StatusBadRequest)
			return
		}

		contracts, err := s.ListAvailableContracts(ctx, filters)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error listing contracts: %v", err), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(contracts); err != nil {
			http.Error(w, "JSON Encoding error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
	}
}

type GuildContractLister interface {
	ListGuildContracts(ctx context.Context, guildID uuid.UUID, filter SearchFilters) ([]ListContractsResponse, error)
}

func ListGuildContracts(s GuildContractLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var filters SearchFilters
		if err := json.NewDecoder(r.Body).Decode(&filters); err != nil {
			http.Error(w, fmt.Sprintf("Bad Request: %v", err), http.StatusBadRequest)
			return
		}

		guildID, ok := auth.GuildIDFromContext(ctx)
		if !ok {
			http.Error(w, "Missing Guild ID Claim", http.StatusUnauthorized)
			return
		}

		contracts, err := s.ListGuildContracts(ctx, guildID, filters)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error listing contracts: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(contracts); err != nil {
			http.Error(w, "JSON Encoding error", http.StatusInternalServerError)
			return
		}
	}
}

type AvailableContractDescriber interface {
	GetAvailableContractDetails(ctx context.Context, contractID uuid.UUID) (ContractDetailsResponse, error)
}

func GetAvailableContractDetails(a AvailableContractDescriber) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		idStr := chi.URLParam(r, "id")
		contract_id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, "Invalid contract ID", http.StatusBadRequest)
			return
		}

		contract, err := a.GetAvailableContractDetails(ctx, contract_id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Request Failed: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(contract); err != nil {
			http.Error(w, "JSON Encode error", http.StatusInternalServerError)
			return
		}
	}
}

type ActiveContractDescriber interface {
	GetActiveContractDetails(ctx context.Context, contractID uuid.UUID) (ActiveContractDetailsResponse, error)
}

func GetActiveContractDetails(a ActiveContractDescriber) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		idStr := chi.URLParam(r, "id")
		contract_id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, "Invalid contract ID", http.StatusBadRequest)
			return
		}

		contract, err := a.GetActiveContractDetails(ctx, contract_id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Request Failed: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(contract); err != nil {
			http.Error(w, "JSON Encode error", http.StatusInternalServerError)
			return
		}
	}
}

type ContractUpdater interface {
	SetContractStatus(ctx context.Context, cs SetContractStatusRequest) error
}

func UpdateContract(c ContractUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var request SetContractStatusRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, fmt.Sprintf("Bad Request: %v", err), http.StatusBadRequest)
			return
		}

		if err := c.SetContractStatus(ctx, request); err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
