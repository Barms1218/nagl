package adventurers

import (
	"github.com/go-chi/chi/v5"
)

func Routes(service *AdventurerService) chi.Router {
	r := chi.NewRouter()

	r.Get("/recruitable", ListRecruitableAdventurers(service))
	r.Get("/details/{id}", GetDetails(service))
	r.Get("/guild", ListGuildMembers(service))
	r.Get("/salary/{id}", GetUpkeepCost(service))
	return r
}
