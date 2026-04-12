package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/maximrozinkevich/daylik/pkg/config"
	"github.com/maximrozinkevich/daylik/pkg/logger"
	"github.com/maximrozinkevich/daylik/pkg/postgres"

	jwt_adapter "github.com/maximrozinkevich/daylik/internal/adapters/jwt"
	deliveryhttp "github.com/maximrozinkevich/daylik/internal/delivery/http"
	habithandler "github.com/maximrozinkevich/daylik/internal/delivery/http/habit"
	userhandler "github.com/maximrozinkevich/daylik/internal/delivery/http/user"
	habitrepo "github.com/maximrozinkevich/daylik/internal/repository/postgres/habit"
	refreshtokenrepo "github.com/maximrozinkevich/daylik/internal/repository/postgres/refresh_token"
	userrepo "github.com/maximrozinkevich/daylik/internal/repository/postgres/user"
	habitusecase "github.com/maximrozinkevich/daylik/internal/usecase/habit"
	"github.com/maximrozinkevich/daylik/internal/usecase/user"
)

func main() {
	log := logger.New()

	cfg, err := config.Load()
	if err != nil {
		log.Error("load config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := postgres.New(ctx, cfg.Postgres.DSN())
	if err != nil {
		log.Error("connect to postgres", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pool.Close()

	txm := postgres.NewTxManager(pool)

	// Adapters
	tokens := jwt_adapter.New(cfg.JWT.Secret, cfg.JWT.AccessTTL)

	// Repositories
	userRepo := userrepo.New(pool)
	refreshTokenRepo := refreshtokenrepo.New(pool)
	habitRepo := habitrepo.New(pool)

	// Services
	userSrv := user.New(userRepo, refreshTokenRepo, tokens, txm, cfg.JWT.RefreshTTL)
	habitSrv := habitusecase.New(habitRepo)

	// Handlers
	userHandler := userhandler.NewHandler(userSrv)
	habitHandler := habithandler.NewHandler(habitSrv)

	// Router
	router := deliveryhttp.NewRouter(log, tokens, userHandler, habitHandler)

	addr := fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	errCh := make(chan error, 1)
	go func() {
		log.Info("server starting", slog.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case <-quit:
		log.Info("shutting down server")
	case err := <-errCh:
		log.Error("server failed", slog.String("error", err.Error()))
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("graceful shutdown failed", slog.String("error", err.Error()))
	}

	log.Info("server stopped")
}
