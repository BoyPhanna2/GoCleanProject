package database

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"

	"go.uber.org/fx"
	"go.uber.org/zap"
	_ "modernc.org/sqlite"

	"myapp/internal/config"
)

func NewSQLiteDB(lc fx.Lifecycle, logger *zap.Logger, cfg *config.Config) (*sql.DB, error) {
	dbPath := cfg.DBPath
	// Ensure directory exists if there's any
	dir := filepath.Dir(dbPath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// SQLite single-writer constraint
	db.SetMaxOpenConns(1)

	// Hook into Fx lifecycle
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("pinging database", zap.String("path", dbPath))
			return db.PingContext(ctx)
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("closing database")
			return db.Close()
		},
	})

	return db, nil
}

// Module provides the sqlite DB and automatically runs migrations.
// We'll use a struct to inject configuration for dbPath if we want, but for now we'll rely on the Config module.
var Module = fx.Provide(NewSQLiteDB)
