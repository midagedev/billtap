package storage

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"
)

var retainedSQLiteTables = map[string]bool{
	"schema_migrations": true,
	"runtime_metadata":  true,
}

// SQLiteTableCounts returns row counts for user-data tables in a SQLite-backed
// store. Migration and runtime metadata are intentionally omitted.
func SQLiteTableCounts(ctx context.Context, store Store) (map[string]int, error) {
	db, err := sqliteDB(store)
	if err != nil {
		return nil, err
	}
	tables, err := sqliteUserTables(ctx, db)
	if err != nil {
		return nil, err
	}
	counts := make(map[string]int, len(tables))
	for _, table := range tables {
		var count int
		if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+quoteSQLiteIdentifier(table)).Scan(&count); err != nil {
			return nil, fmt.Errorf("count %s: %w", table, err)
		}
		counts[table] = count
	}
	return counts, nil
}

// ResetSQLiteData deletes all persisted user data from a SQLite-backed store.
// It keeps schema_migrations and runtime_metadata so the database remains ready
// for immediate reuse.
func ResetSQLiteData(ctx context.Context, store Store) error {
	db, err := sqliteDB(store)
	if err != nil {
		return err
	}
	tables, err := sqliteUserTables(ctx, db)
	if err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, "PRAGMA foreign_keys = OFF"); err != nil {
		return err
	}
	defer func() {
		_, _ = db.ExecContext(context.Background(), "PRAGMA foreign_keys = ON")
	}()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	for _, table := range tables {
		if _, err := tx.ExecContext(ctx, "DELETE FROM "+quoteSQLiteIdentifier(table)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("delete %s: %w", table, err)
		}
	}
	return tx.Commit()
}

func sqliteDB(store Store) (*sql.DB, error) {
	if store == nil {
		return nil, fmt.Errorf("sqlite store is not open")
	}
	withDB, ok := store.(interface{ DB() *sql.DB })
	if !ok || withDB.DB() == nil {
		return nil, fmt.Errorf("storage backend is not sqlite-backed")
	}
	return withDB.DB(), nil
}

func sqliteUserTables(ctx context.Context, db *sql.DB) ([]string, error) {
	rows, err := db.QueryContext(ctx, `SELECT name FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%' ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		if retainedSQLiteTables[name] {
			continue
		}
		tables = append(tables, name)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	sort.Strings(tables)
	return tables, nil
}

func quoteSQLiteIdentifier(value string) string {
	return `"` + strings.ReplaceAll(value, `"`, `""`) + `"`
}
