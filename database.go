package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

func initDB() *sql.DB {
	Database, err := sql.Open("sqlite3", filepath.Join(PodDirectory, "data.db"))
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
	pubDate integer NOT NULL,
	duration NOT NULL,
	size NOT NULL
	);
	CREATE INDEX IF NOT EXISTS addedDate_idx on episodes (addedDate);
	CREATE INDEX IF NOT EXISTS UUID_idx on episodes (UUID);
	`

	_, err = Database.Exec(initcommand)
	if err != nil {
		panic(err)
	}

	init1_0_1 := `CREATE TABLE IF NOT EXISTS meta
	(rowid integer NOT NULL PRIMARY KEY AUTOINCREMENT,
	name text NOT NULL UNIQUE,
	value text NOT NULL);
	CREATE INDEX IF NOT EXISTS metaName_idx on meta (name)`

	_, err = Database.Exec(init1_0_1)
	if err != nil {
		panic(err)
	}

	versionInit := `INSERT OR IGNORE INTO meta (name, value) VALUES ("version", "1.0.1");`

	_, err = Database.Exec(versionInit)
	if err != nil {
		panic(err)
	}

	versionSet := `INSERT INTO meta (name, value) VALUES ("version", "1.0.1")
					ON CONFLICT(name) DO UPDATE SET
					value=excluded.value;`

	_, err = Database.Exec(versionSet)
	if err != nil {
		panic(err)
	}

	return Database
}
