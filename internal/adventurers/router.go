package adventurers

import (
	"github.com/go-chi/chi/v5"
)

func Routes(service *AdventurerService) chi.Router {
	r := chi.NewRouter()

	r.Mount("/recruitable", ListRecruitableAdventurers(service))
	r.Mount("/details/{id}", GetDetails(service))
	r.Mount("/guild", ListGuildMembers(service))
	return r
}
