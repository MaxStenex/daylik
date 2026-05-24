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
	delivery_http "github.com/maximrozinkevich/daylik/internal/delivery/http"
	habit_handler "github.com/maximrozinkevich/daylik/internal/delivery/http/habit"
	habits_log_handler "github.com/maximrozinkevich/daylik/internal/delivery/http/habit_log"
	user_handler "github.com/maximrozinkevich/daylik/internal/delivery/http/user"
	habit_repo "github.com/maximrozinkevich/daylik/internal/repository/postgres/habit"
	habit_log_repo "github.com/maximrozinkevich/daylik/internal/repository/postgres/habit_log"
	refresh_token_repo "github.com/maximrozinkevich/daylik/internal/repository/postgres/refresh_token"
	user_repo "github.com/maximrozinkevich/daylik/internal/repository/postgres/user"
	habit_usecase "github.com/maximrozinkevich/daylik/internal/usecase/habit"
	habit_log_usecase "github.com/maximrozinkevich/daylik/internal/usecase/habit_log"
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
	userRepo := user_repo.New(pool)
	refreshTokenRepo := refresh_token_repo.New(pool)
	habitRepo := habit_repo.New(pool)
	habitLogRepo := habit_log_repo.New(pool)

	// Services
	userSrv := user.New(userRepo, refreshTokenRepo, tokens, txm, cfg.JWT.RefreshTTL)
	habitSrv := habit_usecase.New(habitRepo)
	habitLogSrv := habit_log_usecase.New(habitLogRepo, habitRepo)

	// Handlers
	userHandler := user_handler.New(userSrv, cfg.JWT.RefreshTTL, cfg.HTTP.CookieSecure)
	habitHandler := habit_handler.New(habitSrv)
	habitLogHandler := habits_log_handler.New(habitLogSrv)

	// Router
	router := delivery_http.NewRouter(log, tokens, userHandler, habitHandler, habitLogHandler)

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
