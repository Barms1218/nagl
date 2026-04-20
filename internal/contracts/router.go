package contracts

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Routes(service *ContractService) http.Handler {
	r := chi.NewRouter()
	r.Patch("/claim/{id}", ClaimContract(service))
	r.Patch("/start", StartContract(service))
	r.Get("/guild", ListGuildContracts(service))
	r.Get("/available", ListAvailableContracts(service))
	r.Get("/view_active/{id}", GetActiveContractDetails(service))
	r.Get("/view_available/{id}", GetAvailableContractDetails(service))

	return r
}
