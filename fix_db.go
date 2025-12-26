package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func main() {
	dbPath := "bin/memos/data/memos_dev.db"
	// Use the absolute path if possible or ensure CWD is correct.
	// We are running from memos-0244 root.
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}
	fmt.Println("Connected to DB.")

	// 1. Add missing columns
	columnsToAdd := []struct {
		Name string
		SQL  string
	}{
		{"type", "ALTER TABLE tickets ADD COLUMN type TEXT NOT NULL DEFAULT 'TASK';"},
		{"tags", "ALTER TABLE tickets ADD COLUMN tags TEXT NOT NULL DEFAULT '[]';"},
	}

	for _, col := range columnsToAdd {
		// Check if exists first to be safe
		rows, err := db.Query("PRAGMA table_info(tickets);")
		if err != nil {
			log.Fatal(err)
		}
		exists := false
		for rows.Next() {
			var cid int
			var name string
			var typeStr string
			var notnull int
			var dfltValue *string
			var pk int
			rows.Scan(&cid, &name, &typeStr, &notnull, &dfltValue, &pk)
			if name == col.Name {
				exists = true
				break
			}
		}
		rows.Close()

		if !exists {
			fmt.Printf("Adding column %s...\n", col.Name)
			if _, err := db.Exec(col.SQL); err != nil {
				log.Fatalf("Failed to execute SQL for %s: %v", col.Name, err)
			}
			fmt.Printf("Success: %s added.\n", col.Name)
		} else {
			fmt.Printf("Column %s already exists. Skipping.\n", col.Name)
		}
	}

	// 2. Update migration history to 0.25.2
	// Check if 0.25.2 exists
	rows, err := db.Query("SELECT version FROM migration_history WHERE version = '0.25.2'")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	if !rows.Next() {
		fmt.Println("Inserting version 0.25.2 into migration_history...")
		if _, err := db.Exec("INSERT INTO migration_history (version) VALUES ('0.25.2')"); err != nil {
			log.Fatalf("Failed to update history: %v", err)
		}
		fmt.Println("History updated.")

		// Also insert 0.25.1 just for completeness if not exists
		if _, err := db.Exec("INSERT OR IGNORE INTO migration_history (version) VALUES ('0.25.1')"); err != nil {
			// Ignore error here as it's optional
			fmt.Println("Warning: failed to insert 0.25.1 (might be fine):", err)
		}
	} else {
		fmt.Println("Version 0.25.2 already in history.")
	}

	fmt.Println("Database fix complete.")
}
