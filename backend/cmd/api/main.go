package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/martin-aziz/scopra/backend/src/api"
	"github.com/martin-aziz/scopra/backend/src/config"
	"github.com/martin-aziz/scopra/backend/src/repositories"
	"github.com/martin-aziz/scopra/backend/src/server"
	"github.com/martin-aziz/scopra/backend/src/services"
	"github.com/martin-aziz/scopra/backend/src/utils"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	logger := utils.NewLogger("nexus-mcp-gateway")
	ctx := context.Background()

	dbPool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("database connection failed", map[string]any{"error": err.Error()})
		os.Exit(1)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(ctx); err != nil {
		logger.Error("database ping failed", map[string]any{"error": err.Error()})
		os.Exit(1)
	}

	redisOptions, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		logger.Error("invalid redis url", map[string]any{"error": err.Error()})
		os.Exit(1)
	}
	cache := redis.NewClient(redisOptions)
	defer cache.Close()

	if err := cache.Ping(ctx).Err(); err != nil {
		logger.Error("redis ping failed", map[string]any{"error": err.Error()})
		os.Exit(1)
	}

	repository := repositories.NewPostgresRepository(dbPool)
	tokenService := services.NewTokenService(cfg.JWTSecret, cfg.JWTIssuer, cfg.JWTAudience, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)
	authService := services.NewAuthService(repository, tokenService, services.BcryptHasher{})
	connectorRegistry := services.NewConnectorRegistry()
	toolCallService := services.NewToolCallService(repository, connectorRegistry, cfg.ApprovalMode)
	handler := api.NewHandler(authService, toolCallService, repository, dbPool, cache, logger)
	app := server.New(cfg, handler, tokenService)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-shutdown
		logger.Info("shutdown signal received", map[string]any{})
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = app.ShutdownWithContext(ctx)
	}()

	listenAddress := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	logger.Info("gateway listening", map[string]any{"address": listenAddress})
	if err := app.Listen(listenAddress); err != nil {
		logger.Error("gateway stopped", map[string]any{"error": err.Error()})
		os.Exit(1)
	}
}
