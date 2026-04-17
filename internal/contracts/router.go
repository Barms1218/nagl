package contracts

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Routes(service *ContractService) http.Handler {
	r := chi.NewRouter()

	return r
}
