package app

import (
	"crypto/ecdsa"

	"github.com/Barms1218/nagl/internal/adventurers"
	"github.com/Barms1218/nagl/internal/auth"
	"github.com/Barms1218/nagl/internal/contracts"
	"github.com/Barms1218/nagl/internal/guild"
	"github.com/Barms1218/nagl/internal/procedural"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type App struct {
	guildService      *guild.GuildService
	proceduralService *procedural.ProceduralService
	contractService   *contracts.ContractService
	adventurerService *adventurers.AdventurerService
	privateKey        *ecdsa.PrivateKey
}

func NewApp(
	gs *guild.GuildService,
	ps *procedural.ProceduralService,
	cs *contracts.ContractService,
	as *adventurers.AdventurerService,
	pk *ecdsa.PrivateKey) *App {
	return &App{
		guildService:      gs,
		proceduralService: ps,
		contractService:   cs,
		adventurerService: as,
		privateKey:        pk,
	}
}

func (a *App) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Mount("/guilds", guild.Routes(a.guildService, a.privateKey))

	r.Group(func(r chi.Router) {
		r.Use(auth.JWTMiddleware(&a.privateKey.PublicKey))
		r.Mount("/contracts", contracts.Routes(a.contractService))
		r.Mount("/generate", procedural.Routes(a.proceduralService))
		r.Mount("adventurers", adventurers.Routes(a.adventurerService))
	})

	return r
}
