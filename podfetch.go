package main

import (
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// getPod : Gets the podcast and its metadata and tries to write
// that to the database
func getPod(videoURL string) (episode, error) {
	epuuid := uuid.New()
	uuidbin, _ := epuuid.MarshalBinary()
	uuidstring := hex.EncodeToString(uuidbin)

	fullpath := filepath.Join(PodDirectory, uuidstring)
	jsonBlob, fileSize, err := getAndConvert(videoURL, fullpath)

	if err != nil {
		// TODO: clean up files when this operation fails
		return episode{}, err
	}

	type youtubeData struct {
		Title       string
		Description string
		UploadDate  string `json:"upload_date"`
		Duration    int64  `json:"duration"`
		WebpageURL  string `json:"webpage_url"`
	}

	var ytd youtubeData
	json.Unmarshal(jsonBlob, &ytd)

	var fetchedEp episode

	fetchedEp.title = ytd.Title
	fetchedEp.description = ytd.Description
	fetchedEp.addedDate = time.Now()
	fetchedEp.pubDate, err = time.Parse("20060102", ytd.UploadDate)
	if err != nil {
		fetchedEp.pubDate = time.Unix(0, 0)
	}
	fetchedEp.duration = time.Duration(ytd.Duration) * time.Second

	fetchedEp.URL = ytd.WebpageURL
	fetchedEp.UUID = epuuid
	fetchedEp.size = fileSize

	_, err = Database.Exec(`INSERT INTO episodes(
		UUID,
		title,
		description,
		URL,
		addedDate,
		pubDate,
		duration,
		size
		)
		VALUES(?, ?, ?, ?, ?, ?, ?, ?);
		`, fetchedEp.UUID, fetchedEp.title, fetchedEp.description,
		fetchedEp.URL, fetchedEp.addedDate.Unix(), fetchedEp.pubDate.Unix(),
		int64(fetchedEp.duration/time.Second), fetchedEp.size)

	if err != nil {
		return fetchedEp, err
	}

	u, err := url.Parse(Config.BaseURL)
	if err != nil {
		panic("Configured Base URL is invalid")
	}

	u.Path = path.Join("/instapod/feed/")

	resp, err := http.PostForm("https://overcast.fm/ping",
		url.Values{"urlprefix": {u.String()}})

	if err != nil {
		log.Println(resp)
		log.Println(err)
	}

	return fetchedEp, nil
}

// This is basically a wrapper around youtube-dl, taking a URL
// and filename (without an extension) and returning the video
// metadata (or error), plus placing the audio in an mp3 file
// at the passed path
func getAndConvert(videoURL string, name string) ([]byte, int64, error) {
	command := exec.Command("youtube-dl",
		"-x", "--audio-format", "mp3",
		// E(-x)tract audio in mp3 format
		"--audio-quality", "64k", "--embed-thumbnail",
		// 64kbps baby
		"--add-metadata", "--print-json",
		// Embed metadata in mp3, plus print metadata to stdio
		"--postprocessor-args", "-ac 1",
		// FFmpeg flags to mix down to mono (nb not to double-quote -ac 1)
		"-o", name+".%(ext)s",
		// If you pass the extension directly youtube-dl acts up so you use a
		// placeholder
		// "-s",
		// -s skips downloading the video to save time debugging
		videoURL)
	output, err := command.Output()

	if err != nil {
		panic(err)
	}

	stat, err := os.Stat(name + ".mp3")

	if err != nil {
		panic(err)
	}

	return output, stat.Size(), nil
	// return output, nil, 0
}

func canGetVideo(videoURL string) ([]byte, error) {
	command := exec.Command("youtube-dl", "-j",
		videoURL)

	return command.Output()
}
