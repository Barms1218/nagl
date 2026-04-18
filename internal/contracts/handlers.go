package contracts

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type ContractClaimer interface {
	ClaimContract(ctx context.Context, c ContractClaimRequest) error
}

func ClaimContract(s ContractClaimer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var claimRequest ContractClaimRequest
		if err := json.NewDecoder(r.Body).Decode(&claimRequest); err != nil {
			http.Error(w, "Invalid Request Body", http.StatusBadRequest)
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
			http.Error(w, "Bad Request Body: %w", http.StatusBadRequest)
			return
		}

		if err := s.StartContract(ctx, request); err != nil {
			http.Error(w, "Contract Start Failed: %w", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

type ContractLister interface {
	ListContracts(ctx context.Context, filter SearchFilters) ([]ListContractsResponse, error)
}

func ListContracts(s ContractLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var filters SearchFilters
		if err := json.NewDecoder(r.Body).Decode(&filters); err != nil {
			http.Error(w, "Bad Request: %w", http.StatusBadRequest)
			return
		}

		contracts, err := s.ListContracts(ctx, filters)
		if err != nil {
			http.Error(w, "Error listing contracts: %w", http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(contracts); err != nil {
			http.Error(w, "JSON Encoding error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
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
			http.Error(w, "Bad Request: %w", http.StatusBadRequest)
			return
		}

		if err := c.SetContractStatus(ctx, request); err != nil {
			http.Error(w, "%w", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
