package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/maximrozinkevich/daylik/pkg/config"
	"github.com/maximrozinkevich/daylik/pkg/logger"
	"github.com/pressly/goose/v3"
)

const migrationsDir = "migrations"

func main() {
	log := logger.New()
	if err := run(); err != nil {
		log.Error("migrate", slog.String("error", err.Error()))
		os.Exit(1)
	}
	log.Info("migrations applied successfully")
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	db, err := sql.Open("pgx", cfg.Postgres.DSN())
	if err != nil {
		return err
	}
	defer db.Close() //nolint:errcheck

	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		return err
	}

	if err := goose.UpContext(ctx, db, migrationsDir); err != nil {
		return err
	}

	return nil
}
