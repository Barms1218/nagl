package main

import (
	"context"
	"errors"
	"github.com/Barms1218/nagl/internal/adventurers"
	"github.com/Barms1218/nagl/internal/app"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Barms1218/nagl/internal/contracts"
	"github.com/Barms1218/nagl/internal/database"
	"github.com/Barms1218/nagl/internal/guild"
	"github.com/Barms1218/nagl/internal/procedural"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
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
	a := adventurers.NewAdventurerService(store)

	app := app.NewApp(g, p, c, a, privateKey)

	r := app.Routes()

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP server error", "Error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	shuwdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shuwdownCtx); err != nil {
		slog.Error("Graceful shutdown failed", "Error:", err)
		os.Exit(1)
	}
	log.Println("Closing database connection...")
	log.Println("Shutdown complete. Goodbye!")
}
