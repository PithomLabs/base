package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func main() {
	dbPath := "bin/memos/data/memos_dev.db"
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}

	// Check migration history
	rows, err := db.Query("SELECT version FROM migration_history ORDER BY created_ts DESC LIMIT 1;")
	if err == nil {
		defer rows.Close()
		if rows.Next() {
			var version string
			if err := rows.Scan(&version); err != nil {
				log.Println("Error scanning version:", err)
			} else {
				fmt.Printf("Current DB Version: %s\n", version)
			}
		}
	} else {
		fmt.Println("Error querying migration_history (might not exist):", err)
	}

	// Check tickets columns
	rows, err = db.Query("PRAGMA table_info(tickets);")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Table 'tickets' columns:")
	hasType := false
	for rows.Next() {
		var cid int
		var name string
		var typeStr string
		var notnull int
		var dfltValue *string
		var pk int
		rows.Scan(&cid, &name, &typeStr, &notnull, &dfltValue, &pk)
		fmt.Printf(" - %s (%s)\n", name, typeStr)
		if name == "type" {
			hasType = true
		}
	}

	if !hasType {
		fmt.Println("MISSING 'type' column!")
	} else {
		fmt.Println("'type' column exists.")
	}
}
