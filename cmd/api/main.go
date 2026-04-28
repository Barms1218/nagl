package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Barms1218/nagl/internal/adventurers"
	"github.com/Barms1218/nagl/internal/app"
	"github.com/Barms1218/nagl/internal/contracts"
	"github.com/Barms1218/nagl/internal/database"
	"github.com/Barms1218/nagl/internal/guild"
	"github.com/Barms1218/nagl/internal/procedural"
	"github.com/Barms1218/nagl/internal/workers"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
)

func main() {
	bgCtx := context.Background()
	dbURL := os.Getenv("DB_URL")
	dbConn, err := pgxpool.New(bgCtx, dbURL)
	if err != nil {
		log.Fatal("Database connection failed.")
	}

	store := database.NewStore(dbConn)
	client := anthropic.NewClient(
		option.WithAPIKey(os.Getenv("ANTHROPIC_API_KEY")),
	)
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URL"),
	})

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
	contracts := contracts.NewContractService(rdb, store)
	procedural := procedural.NewProceduralService(&client, store)
	adventurers := adventurers.NewAdventurerService(store)
	workers := workers.NewWorkerService(rdb, store, contracts)

	app := app.NewApp(
		logger,
		guilds,
		procedural,
		contracts,
		adventurers,
		workers,
		privateKey,
	)

	c.AddFunc("@every 4h", func() {
		app.Logger.Info("Starting Procedural Generation", "Type", "Adventurer")
		if err := app.ProceduralService.GenerateAdventurer(bgCtx); err != nil {
			app.Logger.Error("Adventurer Generation Failed", "Error", err)
		}
	})

	c.AddFunc("@every 1h30m", func() {
		app.Logger.Info("Starting Procedural Generation", "Type", "Contract")
		if err := app.ProceduralService.GenerateContract(bgCtx); err != nil {
			app.Logger.Error("Contract Generation Failed", "Error", err)
		}
	})

	c.AddFunc("@every 1h", func() {
		app.Logger.Info("Checking expired contracts")
		if err := app.ContractService.CheckExpiredContracts(bgCtx); err != nil {
			app.Logger.Error("Contract Resolution Failed", "Error", err)
		}
	})

	c.Start()
	r := app.Routes()

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Get("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if err := dbConn.Ping(r.Context()); err != nil {
			logger.Error("Ready Check Failed", "Error", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

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
