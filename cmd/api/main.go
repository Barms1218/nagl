package main

import (
	"context"
	"net/http"
	"os"

	"github.com/Barms1218/nagl/internal/app"

	"log"

	"github.com/Barms1218/nagl/internal/contracts"
	"github.com/Barms1218/nagl/internal/database"
	"github.com/Barms1218/nagl/internal/guild"
	"github.com/Barms1218/nagl/internal/procedural"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

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
	p := procedural.NewProceduralService(&client, store)

	app := app.NewApp(g, p, c, privateKey)

	r := app.Routes()

	http.ListenAndServe(":3000", r)
}
