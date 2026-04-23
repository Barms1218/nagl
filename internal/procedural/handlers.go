package procedural

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type PartyGenerator interface {
	GenerateParty(ctx context.Context, r GeneratePartyRequest) (GeneratedParty, error)
}

func RequestParty(p PartyGenerator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var request GeneratePartyRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Bad Request Body", http.StatusBadRequest)
			return
		}

		party, err := p.GenerateParty(ctx, request)
		if err != nil {
			http.Error(w, fmt.Sprintf("Party request failed: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(party); err != nil {
			http.Error(w, "JSON Encode error", http.StatusInternalServerError)
			return
		}
	}
}
