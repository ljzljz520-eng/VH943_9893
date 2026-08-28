package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"campgear/internal/domain"
	_ "modernc.org/sqlite"
)

var ErrNotFound = domain.ErrItemNotFound

type Repository struct {
	db *sql.DB
}

func Open(path string) (*Repository, error) {
	if path == "" {
		return nil, fmt.Errorf("database path is required")
	}
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	repo := &Repository{db: db}
	if err := repo.configure(context.Background()); err != nil {
		db.Close()
		return nil, err
	}
	if err := repo.migrate(context.Background()); err != nil {
		db.Close()
		return nil, err
	}
	return repo, nil
}

func OpenMemory() (*Repository, error) {
	return Open("file:campgear-memory?mode=memory&cache=shared")
}

func (r *Repository) DB() *sql.DB { return r.db }

func (r *Repository) Close() error {
	if r == nil || r.db == nil {
		return nil
	}
	return r.db.Close()
}

func (r *Repository) configure(ctx context.Context) error {
	for _, statement := range []string{"PRAGMA foreign_keys = ON", "PRAGMA busy_timeout = 5000"} {
		if _, err := r.db.ExecContext(ctx, statement); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) migrate(ctx context.Context) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS inventory_items (id TEXT PRIMARY KEY, sku TEXT NOT NULL UNIQUE, name TEXT NOT NULL, category TEXT NOT NULL, daily_rate INTEGER NOT NULL, deposit INTEGER NOT NULL, stock INTEGER NOT NULL, available INTEGER NOT NULL, maintenance_status TEXT NOT NULL, storage_bin TEXT NOT NULL, listed INTEGER NOT NULL, version INTEGER NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS rental_records (id TEXT PRIMARY KEY, customer TEXT NOT NULL, start_date TEXT NOT NULL, end_date TEXT NOT NULL, status TEXT NOT NULL, total INTEGER NOT NULL, deposit_held INTEGER NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS rental_lines (rental_id TEXT NOT NULL REFERENCES rental_records(id) ON DELETE CASCADE, item_id TEXT NOT NULL REFERENCES inventory_items(id), quantity INTEGER NOT NULL, days INTEGER NOT NULL, rate INTEGER NOT NULL, deposit INTEGER NOT NULL, subtotal INTEGER NOT NULL, PRIMARY KEY(rental_id, item_id))`,
		`CREATE TABLE IF NOT EXISTS maintenance_orders (id TEXT PRIMARY KEY, item_id TEXT NOT NULL REFERENCES inventory_items(id), reason TEXT NOT NULL, opened_date TEXT NOT NULL, closed_date TEXT NOT NULL, status TEXT NOT NULL, technician TEXT NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS audit_events (id TEXT PRIMARY KEY, entity_type TEXT NOT NULL, entity_id TEXT NOT NULL, action TEXT NOT NULL, actor TEXT NOT NULL, occurred_at TEXT NOT NULL, details TEXT NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS staff_members (id TEXT PRIMARY KEY, name TEXT NOT NULL, role TEXT NOT NULL, active INTEGER NOT NULL)`,
	}
	for _, statement := range statements {
		if _, err := r.db.ExecContext(ctx, statement); err != nil {
			return err
		}
	}
	return nil
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func intBool(value int64) bool { return value != 0 }

func wrapNotFound(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return err
}

func normalizeQuery(value string) string { return strings.TrimSpace(value) }
