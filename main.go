package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// PodDirectory : Directory of instapod's data
var PodDirectory string

// Database : Common SQLite database for all episodes
var Database *sql.DB

type episode struct {
	title       string
	description string
	URL         string
	UUID        uuid.UUID
	addedDate   time.Time
	pubDate     time.Time
}

func main() {
	// Initialization
	checkDeps()
	PodDirectory = getDir()
	Database = initDB()

	defer Database.Close()

	if len(os.Args) < 2 {
		fmt.Println("Please supply a URL as an argument")
		os.Exit(1)
	}

	url := os.Args[1]

	getPod(url)
}
