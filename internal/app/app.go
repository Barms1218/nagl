package app

import (
	"crypto/ecdsa"
	"log/slog"

	"github.com/Barms1218/nagl/internal/adventurers"
	"github.com/Barms1218/nagl/internal/auth"
	"github.com/Barms1218/nagl/internal/contracts"
	"github.com/Barms1218/nagl/internal/guild"
	"github.com/Barms1218/nagl/internal/procedural"
	"github.com/Barms1218/nagl/internal/workers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type App struct {
	Logger            *slog.Logger
	GuildService      *guild.GuildService
	ProceduralService *procedural.ProceduralService
	ContractService   *contracts.ContractService
	AdventurerService *adventurers.AdventurerService
	WorkerService     *workers.WorkerService
	privateKey        *ecdsa.PrivateKey
}

func NewApp(
	logger *slog.Logger,
	gs *guild.GuildService,
	ps *procedural.ProceduralService,
	cs *contracts.ContractService,
	as *adventurers.AdventurerService,
	ws *workers.WorkerService,
	pk *ecdsa.PrivateKey) *App {
	return &App{
		Logger:            logger,
		GuildService:      gs,
		ProceduralService: ps,
		ContractService:   cs,
		AdventurerService: as,
		WorkerService:     ws,
		privateKey:        pk,
	}
}

func (a *App) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Mount("/guilds", guild.Routes(a.GuildService, a.privateKey))

	r.Group(func(r chi.Router) {
		r.Use(auth.JWTMiddleware(&a.privateKey.PublicKey))
		r.Mount("/contracts", contracts.Routes(a.ContractService))
		r.Mount("/generate", procedural.Routes(a.ProceduralService))
		r.Mount("/adventurers", adventurers.Routes(a.AdventurerService))
	})

	return r
}
