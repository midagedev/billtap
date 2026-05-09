package storage

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	_ "modernc.org/sqlite"
)

const (
	DriverSQLite = "sqlite"
	DriverMemory = "memory"
)

//go:embed migrations/*.sql
var migrations embed.FS

type Store interface {
	Ping(context.Context) error
	MigrationVersions(context.Context) ([]int, error)
	Close() error
}

type Options struct {
	Driver string
	DSN    string
}

func Open(ctx context.Context, opts Options) (Store, error) {
	switch opts.Driver {
	case "", DriverSQLite:
		return OpenSQLite(ctx, opts.DSN)
	case DriverMemory:
		return NewMemoryStore(), nil
	default:
		return nil, fmt.Errorf("unsupported storage driver %q", opts.Driver)
	}
}

type SQLiteStore struct {
	db *sql.DB
}

func (s *SQLiteStore) DB() *sql.DB {
	if s == nil {
		return nil
	}
	return s.db
}

func OpenSQLite(ctx context.Context, dsn string) (*SQLiteStore, error) {
	if dsn == "" {
		dsn = ".billtap/billtap.db"
	}
	if err := ensureParentDir(dsn); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	db.SetMaxOpenConns(1)

	store := &SQLiteStore{db: db}
	if _, err := db.ExecContext(ctx, "PRAGMA foreign_keys = ON"); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("enable sqlite foreign keys: %w", err)
	}
	if err := runMigrations(ctx, db, migrations, "migrations"); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := store.Ping(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

func (s *SQLiteStore) Ping(ctx context.Context) error {
	if s == nil || s.db == nil {
		return errors.New("sqlite store is not open")
	}
	return s.db.PingContext(ctx)
}

func (s *SQLiteStore) MigrationVersions(ctx context.Context) ([]int, error) {
	if s == nil || s.db == nil {
		return nil, errors.New("sqlite store is not open")
	}
	rows, err := s.db.QueryContext(ctx, "SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []int
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}
	return versions, rows.Err()
}

func (s *SQLiteStore) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

type MemoryStore struct {
	closed atomic.Bool
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (s *MemoryStore) Ping(context.Context) error {
	if s == nil {
		return errors.New("memory store is not open")
	}
	if s.closed.Load() {
		return errors.New("memory store is closed")
	}
	return nil
}

func (s *MemoryStore) MigrationVersions(ctx context.Context) ([]int, error) {
	if s == nil {
		return nil, errors.New("memory store is not open")
	}
	return []int{}, s.Ping(ctx)
}

func (s *MemoryStore) Close() error {
	if s != nil {
		s.closed.Store(true)
	}
	return nil
}

func runMigrations(ctx context.Context, db *sql.DB, migrationFS fs.FS, root string) error {
	if _, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		applied_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`); err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}

	entries, err := fs.ReadDir(migrationFS, root)
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		version, err := migrationVersion(entry.Name())
		if err != nil {
			return err
		}

		var exists bool
		row := db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = ?)", version)
		if err := row.Scan(&exists); err != nil {
			return fmt.Errorf("check migration %s: %w", entry.Name(), err)
		}
		if exists {
			continue
		}

		body, err := fs.ReadFile(migrationFS, filepath.Join(root, entry.Name()))
		if err != nil {
			return fmt.Errorf("read migration %s: %w", entry.Name(), err)
		}

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("begin migration %s: %w", entry.Name(), err)
		}
		if _, err := tx.ExecContext(ctx, string(body)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("execute migration %s: %w", entry.Name(), err)
		}
		if _, err := tx.ExecContext(ctx, "INSERT INTO schema_migrations (version, name, applied_at) VALUES (?, ?, ?)", version, entry.Name(), time.Now().UTC().Format(time.RFC3339)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("record migration %s: %w", entry.Name(), err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %s: %w", entry.Name(), err)
		}
	}
	return nil
}

func migrationVersion(name string) (int, error) {
	prefix, _, found := strings.Cut(name, "_")
	if !found {
		return 0, fmt.Errorf("migration %q must start with a numeric prefix and underscore", name)
	}
	version, err := strconv.Atoi(prefix)
	if err != nil {
		return 0, fmt.Errorf("parse migration version from %q: %w", name, err)
	}
	return version, nil
}

func ensureParentDir(dsn string) error {
	if dsn == ":memory:" || strings.HasPrefix(dsn, "file::memory:") {
		return nil
	}

	path := dsn
	if strings.HasPrefix(path, "file:") {
		path = strings.TrimPrefix(path, "file:")
		if idx := strings.IndexAny(path, "?#"); idx >= 0 {
			path = path[:idx]
		}
	}
	if path == "" || strings.Contains(path, "mode=memory") {
		return nil
	}

	dir := filepath.Dir(path)
	if dir == "." || dir == "" {
		return nil
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create sqlite directory %q: %w", dir, err)
	}
	return nil
}
