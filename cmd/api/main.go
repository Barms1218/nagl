package main

import (
	"context"
	"errors"
	"github.com/Barms1218/nagl/internal/adventurers"
	"github.com/Barms1218/nagl/internal/app"
	"github.com/Barms1218/nagl/internal/contracts"
	"github.com/Barms1218/nagl/internal/database"
	"github.com/Barms1218/nagl/internal/guild"
	"github.com/Barms1218/nagl/internal/procedural"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()
	dbURL := os.Getenv("DB_URL")
	dbConn, err := pgxpool.New(ctx, dbURL)

	store := database.NewStore(dbConn)
	client := anthropic.NewClient(
		option.WithAPIKey(os.Getenv("ANTHROPIC_API_KEY")),
	)
	keyBytes, err := os.ReadFile("ec_private.pem")
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := jwt.ParseECPrivateKeyFromPEM(keyBytes)
	if err != nil {
		log.Fatal(err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	c := cron.New()
	v := validator.New()
	guilds := guild.NewGuildService(store, v, privateKey)
	contracts := contracts.NewContractService(store)
	procedural := procedural.NewProceduralService(&client, store)
	adventurers := adventurers.NewAdventurerService(store)

	app := app.NewApp(
		logger,
		guilds,
		procedural,
		contracts,
		adventurers,
		privateKey,
	)

	c.AddFunc("@every 4h", func() {
		go func() {
			app.Logger.Info("Starting Procedural Generation", "Type", "Adventurer")
			if err := app.ProceduralService.GenerateAdventurer(ctx); err != nil {
				app.Logger.Error("Adventurer Generation Failed", "Error", err)
			}
		}()
	})

	c.AddFunc("@every 1h30m", func() {
		go func() {
			app.Logger.Info("Starting Procedural Generation", "Type", "Contract")
			if err := app.ProceduralService.GenerateContract(ctx); err != nil {
				app.Logger.Error("Contract Generation Failed", "Error", err)
			}
		}()
	})

	c.AddFunc("@every 1h", func() {
		go func() {
			app.Logger.Info("Checking expired contracts")
			if err := app.ContractService.CheckExpiredContracts(ctx); err != nil {
				app.Logger.Error("Contract Resolution Failed", "Error", err)
			}
		}()
	})

	c.Start()
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
