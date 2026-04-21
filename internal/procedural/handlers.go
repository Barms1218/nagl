package procedural

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type AdventurerGenerator interface {
	GenerateAdventurer(ctx context.Context) (GeneratedAdventurer, error)
}

func RequestAdventurer(a AdventurerGenerator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		adventurer, err := a.GenerateAdventurer(ctx)
		if err != nil {
			http.Error(w, "Anthropic API Failure", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(adventurer); err != nil {
			http.Error(w, "JSON Encoding error", http.StatusInternalServerError)
			return
		}

	}
}

type ContractGenerator interface {
	GenerateContract(ctx context.Context) (GeneratedContract, error)
}

func RequestContract(c ContractGenerator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		contract, err := c.GenerateContract(ctx)
		if err != nil {
			http.Error(w, "Anthropic API Failure", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(contract); err != nil {
			http.Error(w, "JSON Encoding error", http.StatusInternalServerError)
			return
		}
	}
}
