package contracts

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Routes(service *ContractService) http.Handler {
	r := chi.NewRouter()
	r.Mount("/claim/{id}", ClaimContract(service))
	r.Mount("/start", StartContract(service))
	r.Mount("/guild", ListGuildContracts(service))
	r.Mount("/available", ListAvailableContracts(service))
	r.Mount("/view_active/{id}", GetActiveContractDetails(service))
	r.Mount("/view_available/{id}", GetAvailableContractDetails(service))

	return r
}
