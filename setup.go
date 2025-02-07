package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/database"
)

func setup(ctx context.Context, db *sql.DB, logger *slog.Logger) error {
	provider, err := goose.NewProvider(database.DialectPostgres, db, os.DirFS("sql/schema"), goose.WithVerbose(true))
	if err != nil {
		return err
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("Context cancelled setup")
		case <-ticker.C:
			if e := provider.Ping(ctx); e == nil {
				break
			} else {
				logger.LogAttrs(ctx, slog.LevelInfo, "provider ping attempt error")
			}
		}
		break
	}

	res, err := provider.Up(ctx)
	if err != nil {
		return err
	}

	for _, mr := range res {
		logger.LogAttrs(ctx, slog.LevelInfo, "migration result", slog.String("result", mr.String()))
	}
	return nil
}
