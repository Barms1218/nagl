package contracts

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

func ClaimContract(s *ContractService) http.HandlerFunc {
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
