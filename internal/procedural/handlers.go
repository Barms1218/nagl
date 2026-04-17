package procedural

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

func RequestAdventurer(p *ProceduralService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		contract, err := p.GenerateAdventurer(ctx)
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
