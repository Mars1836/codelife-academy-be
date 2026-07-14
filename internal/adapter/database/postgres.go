package database

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io/fs"
	"sort"
	"strconv"
	"strings"
	"time"

	migrationfiles "codelife-study-be/migrations"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const migrationLockID int64 = 13820260714

type migration struct {
	version  int64
	name     string
	sql      string
	checksum string
}

func OpenPostgres(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	if databaseURL == "" {
		return nil, nil
	}
	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}
	// Conservative defaults for a 2 GB VPS.
	cfg.MaxConns, cfg.MinConns = 5, 0
	cfg.MaxConnIdleTime, cfg.MaxConnLifetime = 5*time.Minute, 30*time.Minute
	cfg.HealthCheckPeriod = time.Minute
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}

func RunMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	if pool == nil {
		return nil
	}
	if _, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version BIGINT PRIMARY KEY,
			name TEXT,
			checksum TEXT,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`); err != nil {
		return fmt.Errorf("bootstrap schema migrations: %w", err)
	}
	if _, err := pool.Exec(ctx, `ALTER TABLE schema_migrations ADD COLUMN IF NOT EXISTS name TEXT`); err != nil {
		return fmt.Errorf("upgrade schema migrations table: %w", err)
	}
	if _, err := pool.Exec(ctx, `ALTER TABLE schema_migrations ADD COLUMN IF NOT EXISTS checksum TEXT`); err != nil {
		return fmt.Errorf("upgrade schema migrations checksum: %w", err)
	}

	migrations, err := loadMigrations(migrationfiles.Files)
	if err != nil {
		return err
	}
	for _, item := range migrations {
		if err := applyMigration(ctx, pool, item); err != nil {
			return err
		}
	}
	return nil
}

func loadMigrations(files fs.FS) ([]migration, error) {
	entries, err := fs.ReadDir(files, ".")
	if err != nil {
		return nil, fmt.Errorf("read embedded migrations: %w", err)
	}
	items := make([]migration, 0, len(entries))
	seen := make(map[int64]string)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		prefix, _, ok := strings.Cut(entry.Name(), "_")
		if !ok {
			return nil, fmt.Errorf("migration %q must start with a numeric timestamp and underscore", entry.Name())
		}
		version, err := strconv.ParseInt(prefix, 10, 64)
		if err != nil || len(prefix) != 14 {
			return nil, fmt.Errorf("migration %q must use YYYYMMDDHHMMSS_name.sql", entry.Name())
		}
		if _, err := time.Parse("20060102150405", prefix); err != nil {
			return nil, fmt.Errorf("migration %q contains an invalid timestamp", entry.Name())
		}
		if previous, exists := seen[version]; exists {
			return nil, fmt.Errorf("migrations %q and %q have the same timestamp", previous, entry.Name())
		}
		raw, err := fs.ReadFile(files, entry.Name())
		if err != nil {
			return nil, fmt.Errorf("read migration %q: %w", entry.Name(), err)
		}
		if strings.TrimSpace(string(raw)) == "" {
			return nil, fmt.Errorf("migration %q is empty", entry.Name())
		}
		seen[version] = entry.Name()
		checksum := fmt.Sprintf("%x", sha256.Sum256(raw))
		items = append(items, migration{version: version, name: entry.Name(), sql: string(raw), checksum: checksum})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].version < items[j].version })
	return items, nil
}

func applyMigration(ctx context.Context, pool *pgxpool.Pool, item migration) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin migration %s: %w", item.name, err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock($1)`, migrationLockID); err != nil {
		return fmt.Errorf("lock migration %s: %w", item.name, err)
	}
	var storedChecksum *string
	err = tx.QueryRow(ctx, `SELECT checksum FROM schema_migrations WHERE version = $1`, item.version).Scan(&storedChecksum)
	if err == nil {
		if storedChecksum != nil && *storedChecksum != item.checksum {
			return fmt.Errorf("migration %s was modified after being applied", item.name)
		}
		return tx.Commit(ctx)
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("check migration %s: %w", item.name, err)
	}
	if _, err := tx.Exec(ctx, item.sql); err != nil {
		return fmt.Errorf("apply migration %s: %w", item.name, err)
	}
	if _, err := tx.Exec(ctx, `INSERT INTO schema_migrations (version, name, checksum) VALUES ($1, $2, $3)`, item.version, item.name, item.checksum); err != nil {
		return fmt.Errorf("record migration %s: %w", item.name, err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit migration %s: %w", item.name, err)
	}
	return nil
}
