package main

import (
	"context"
	"crypto/ecdsa"
	"os"

	"log"

	"github.com/Barms1218/nagl/internal/contracts"
	"github.com/Barms1218/nagl/internal/database"
	"github.com/Barms1218/nagl/internal/guild"
	"github.com/Barms1218/nagl/internal/procedural"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	guildService      *guild.GuildService
	proceduralService *procedural.ProceduralService
	contractService   *contracts.ContractService
	privateKey        *ecdsa.PrivateKey
}

func NewApp(
	gs *guild.GuildService,
	ps *procedural.ProceduralService,
	cs *contracts.ContractService,
	pk *ecdsa.PrivateKey) *App {
	return &App{
		guildService:      gs,
		proceduralService: ps,
		contractService:   cs,
		privateKey:        pk,
	}
}

func (a *App) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Mount("/guilds", guild.Routes(a.guildService, a.privateKey))
}

func main() {
	ctx := context.Background()
	dbURL := os.Getenv("DB_URL")
	dbConn, err := pgxpool.New(ctx, dbURL)

	//q := database.New(dbConn)
	store := database.NewStore(dbConn)
	client := anthropic.NewClient(
		option.WithAPIKey("my-anthropic-api-key"), // defaults to os.LookupEnv("ANTHROPIC_API_KEY")
	)
	keyBytes, err := os.ReadFile("ec_private.pem")
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := jwt.ParseECPrivateKeyFromPEM(keyBytes)
	if err != nil {
		log.Fatal(err)
	}

	v := validator.New()
	g := guild.NewGuildService(store, v, privateKey)
	c := contracts.NewContractService(store)
	p := procedural.NewProceduralService(client, store)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Mount("/guilds", guild.Routes(g, privateKey))
	r.Mount("/contracts", contracts.Routes(c, privateKey))

}
