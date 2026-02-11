package postgres

import (
	"database/sql"
	"github.com/rh-mithu/rizon/backend/config"
	"log"
	"log/slog"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

func NewMigrator(cfg *config.Config, l *slog.Logger) *bun.DB {
	dsn := cfg.SQLDatabaseURL
	if dsn == "" {
		log.Fatal("DATABASE_URL not found in environment variables")
	}
	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to Postgres: %v", err)
	}

	db := bun.NewDB(sqlDB, pgdialect.New())
	db.WithQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	if err = db.Ping(); err != nil {
		l.Error("failed to ping Postgres", slog.String("error", err.Error()))
		os.Exit(1)
	}
	l.Info("Connected to Postgres successfully")
	return db
}
