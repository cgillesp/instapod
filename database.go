package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

func initDB() *sql.DB {
	db, err := sql.Open("sqlite3", filepath.Join(podDirectory, "data.db"))
	if err != nil {
		fmt.Println("Failed to load ~/.instapod/data.db . Check your permissions")
		os.Exit(1)
	}

	initcommand := `CREATE TABLE IF NOT EXISTS episodes
	(rowid integer NOT NULL PRIMARY KEY AUTOINCREMENT,
	UUID blob NOT NULL UNIQUE,
	title text NOT NULL,
	description text,
	URL text NOT NULL,
	addedDate integer NOT NULL,
	pubDate integer NOT NULL);
	CREATE INDEX IF NOT EXISTS addedDate_idx on episodes (addedDate);
	CREATE INDEX IF NOT EXISTS UUID_idx on episodes (UUID);
	`

	_, err = db.Exec(initcommand)

	if err != nil {
		panic(err)
	}

	return db
}
