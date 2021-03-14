package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

// PodDirectory : Directory of instapod's data
var PodDirectory string

// Database : Common SQLite Database for all episodes
var Database *sql.DB

// Config : Global app configuration
var Config configuration

type episode struct {
	title       string
	description string
	URL         string
	UUID        uuid.UUID
	addedDate   time.Time
	pubDate     time.Time
	duration    time.Duration
	size        int64
}

func main() {
	// Initialization
	checkDeps()

	PodDirectory = getDir()

	logpath := filepath.Join(PodDirectory, "logs.txt")

	logfile, err := os.OpenFile(logpath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logfile)

	Config = getConfig()
	Database = initDB()

	// defer Database.Close()

	serve()
}

func serve() {
	r := mux.NewRouter()

	r.HandleFunc("/instapod/episodes/", addEpisode)
	r.HandleFunc("/instapod/feed/{key}", getFeed)
	r.HandleFunc("/instapod/files/{id}.mp3", getFile)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":19565", nil))
}
