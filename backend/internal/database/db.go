package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	Gorm *gorm.DB
	SQL  *sql.DB
}

func Open(ctx context.Context, log *slog.Logger, path string) (DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return DB{}, fmt.Errorf("mkdir db dir: %w", err)
	}

	gormLogger := logger.New(
		slog.NewLogLogger(log.Handler(), slog.LevelWarn),
		logger.Config{
			SlowThreshold:             500 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	gdb, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return DB{}, fmt.Errorf("open sqlite: %w", err)
	}

	sqldb, err := gdb.DB()
	if err != nil {
		return DB{}, fmt.Errorf("gorm db: %w", err)
	}
	sqldb.SetConnMaxLifetime(30 * time.Minute)
	sqldb.SetMaxOpenConns(10)
	sqldb.SetMaxIdleConns(10)

	if _, err := sqldb.ExecContext(ctx, "PRAGMA foreign_keys = ON;"); err != nil {
		return DB{}, fmt.Errorf("pragma foreign_keys: %w", err)
	}
	if _, err := sqldb.ExecContext(ctx, "PRAGMA journal_mode = WAL;"); err != nil {
		return DB{}, fmt.Errorf("pragma journal_mode: %w", err)
	}
	if _, err := sqldb.ExecContext(ctx, "PRAGMA synchronous = NORMAL;"); err != nil {
		return DB{}, fmt.Errorf("pragma synchronous: %w", err)
	}

	return DB{Gorm: gdb, SQL: sqldb}, nil
}
