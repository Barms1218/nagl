package procedural

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type AdventurerGenerator interface {
	GenerateAdventurer(ctx context.Context) (GeneratedContract, error)
}

func RequestAdventurer(a AdventurerGenerator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		contract, err := a.GenerateAdventurer(ctx)
		if err != nil {
			http.Error(w, "Anthropic API Failure", http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(contract); err != nil {
			http.Error(w, "JSON Encoding error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
	}
}
