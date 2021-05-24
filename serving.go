package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/eduncan911/podcast"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func getFile(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	path := filepath.Join(PodDirectory, id) + ".mp3"
	http.ServeFile(w, r, path)
}

func getEpisodes() []episode {
	rows, err := Database.Query(`SELECT UUID, title, description,
	URL, addedDate, pubDate, duration, size FROM episodes WHERE available=TRUE;`)
	if err != nil {
		panic(err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)

	eps := []episode{}
	for rows.Next() {
		var e episode
		var addedDateInt int64
		var pubDateInt int64
		var durationInt int64
		err := rows.Scan(&e.UUID, &e.title, &e.description,
			&e.URL, &addedDateInt, &pubDateInt, &durationInt, &e.size)

		if err != nil {
			log.Println(err)
			continue
		}

		e.addedDate = time.Unix(addedDateInt, 0)
		e.pubDate = time.Unix(pubDateInt, 0)
		e.duration = time.Duration(durationInt) * time.Second
		eps = append(eps, e)
	}
	return eps
}

func addEpisode(w http.ResponseWriter, r *http.Request) {

	var sentKey string
	var videoURL string

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		sentKey = r.Form.Get("key")
		videoURL = r.Form.Get("url")
	} else if r.Method == http.MethodGet {
		queryStrings := r.URL.Query()
		sentKey = queryStrings.Get("key")
		videoURL = queryStrings.Get("url")
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

	if !safeEquals([]byte(sentKey), Config.addKeyBytes) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	_, err := canGetVideo(videoURL)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := fmt.Fprintf(w, "❌ Failed to Fetch URL: %s", videoURL)
		if err != nil {
			log.Println(err)
		}
		return
	}

	w.WriteHeader(http.StatusAccepted)
	_, err = fmt.Fprintf(w, "✅ Fetching URL: %s", videoURL)
	if err != nil {
		log.Println(err)
	}

	go func() {
		_, err := getPod(videoURL)
		if err != nil {
			log.Println(err)
		}
	}()
}

func getFeed(w http.ResponseWriter, r *http.Request) {
	sentKey := mux.Vars(r)["key"]

	if !safeEquals([]byte(sentKey), Config.readKeyBytes) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write(makePodcastFeed())
	if err != nil {
		log.Println(err)
	}

}

func makePodcastFeed() []byte {
	now := time.Now()
	pod := podcast.New(Config.Title, Config.Link, Config.Description, &now, &now)

	if len(Config.ImageURL) > 0 {
		pod.AddImage(Config.ImageURL)
	}

	eps := getEpisodes()

	for _, e := range eps {

		item := podcast.Item{
			Title:       e.title,
			Description: e.description,
			Link:        e.URL,
			PubDate:     &e.addedDate,
			GUID:        getHexUUID(e.UUID),
		}

		if len(item.Description) == 0 {
			item.Description = "(No description)"
		}
		if len(item.Title) == 0 {
			URL, err := url.Parse(item.Link)
			if err == nil {
				item.Title = fmt.Sprintf("%s, %s", URL.Hostname(),
					e.pubDate.Format("1/2/06"))
			} else {
				item.Title = "Title unavailable"
			}
		}

		// needs a lot of work
		item.AddEnclosure(getURL(e.UUID), podcast.MP3, e.size)

		item.AddDuration(int64(e.duration / time.Second))
		_, err := pod.AddItem(item)

		if err != nil {
			log.Println(err)
		}
	}

	pod.Generator = "Instapod v0.1.0"
	return pod.Bytes()
}
