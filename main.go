package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

// podDirectory : Directory of instapod's data
var podDirectory string

// database : Common SQLite database for all episodes
var database *sql.DB

var addKey string
var readKey string

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
	podDirectory = getDir()
	database = initDB()

	defer database.Close()

	addKey = "hello"
	serve()
}

func serve() {
	r := mux.NewRouter()

	r.HandleFunc("/instapod/episodes", addEpisode).Methods("POST")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":19565", nil))
}

func addEpisode(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	submittedKey := r.Form.Get("key")
	videoURL := r.Form.Get("url")

	if submittedKey != addKey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, err := canGetVideo(videoURL)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "Received. Key: %s \n URL: %s", submittedKey, videoURL)

	go getPod(videoURL)
}
