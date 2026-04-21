package procedural

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Routes(p *ProceduralService) http.Handler {
	r := chi.NewRouter()
	r.Post("/adventurers", RequestAdventurer(p))
	r.Post("/contracts", RequestContract(p))
	r.Post("/parties", RequestParty(p))

	return r
}
