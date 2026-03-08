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

	deliveryhttp "github.com/maximrozinkevich/daylik/internal/delivery/http"
	userhandler "github.com/maximrozinkevich/daylik/internal/delivery/http/user"
	userrepo "github.com/maximrozinkevich/daylik/internal/repository/postgres"
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

	// Repositories
	userRepo := userrepo.NewUserRepository(pool)

	// Services
	userSrv := user.New(userRepo)

	// Handlers
	userHandler := userhandler.NewHandler(userSrv)

	// Router
	router := deliveryhttp.NewRouter(log, userHandler)

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
