package store

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
)

// AddColumnIfNotExists adds a column to a table if it doesn't already exist
// This is a helper for SQLite migrations since ALTER TABLE ADD COLUMN IF NOT EXISTS is not supported
func AddColumnIfNotExists(ctx context.Context, db *sql.DB, tableName, columnName, columnDef string) error {
	// Check if column exists
	query := fmt.Sprintf("PRAGMA table_info(%s)", tableName)
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to query table info: %w", err)
	}
	defer rows.Close()

	exists := false
	for rows.Next() {
		var cid int
		var name, typ string
		var notnull, dfltValue, pk sql.NullString
		if err := rows.Scan(&cid, &name, &typ, &notnull, &dfltValue, &pk); err != nil {
			return fmt.Errorf("failed to scan column info: %w", err)
		}
		if name == columnName {
			exists = true
			break
		}
	}

	if exists {
		slog.Info("Column already exists, skipping", "table", tableName, "column", columnName)
		return nil
	}

	// Add column
	alterSQL := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", tableName, columnName, columnDef)
	slog.Info("Adding column", "table", tableName, "column", columnName, "sql", alterSQL)
	_, err = db.ExecContext(ctx, alterSQL)
	if err != nil {
		return fmt.Errorf("failed to add column: %w", err)
	}

	slog.Info("Column added successfully", "table", tableName, "column", columnName)
	return nil
}

// EnsureTicketBeadsColumns ensures all beads-related columns exist in tickets table
func (s *Store) EnsureTicketBeadsColumns(ctx context.Context) error {
	db := s.driver.GetDB()

	columns := []struct {
		name string
		def  string
	}{
		{"beads_id", "TEXT"},
		{"parent_id", "INTEGER REFERENCES tickets(id)"},
		{"labels", "TEXT DEFAULT '[]'"},
		{"dependencies", "TEXT DEFAULT '[]'"},
		{"discovery_context", "TEXT"},
		{"closed_reason", "TEXT"},
		{"issue_type", "TEXT"},
	}

	for _, col := range columns {
		if err := AddColumnIfNotExists(ctx, db, "tickets", col.name, col.def); err != nil {
			return fmt.Errorf("failed to add column %s: %w", col.name, err)
		}
	}

	// Create unique index on beads_id after column is added
	_, err := db.ExecContext(ctx, "CREATE UNIQUE INDEX IF NOT EXISTS idx_tickets_beads_id ON tickets(beads_id)")
	if err != nil {
		slog.Warn("Failed to create unique index on beads_id, it may already exist", "error", err)
	}

	return nil
}
