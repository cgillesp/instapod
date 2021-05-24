package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// getPod : Gets the podcast and its metadata and tries to write
// that to the database
func getPod(videoURL string) (episode, error) {
	epuuid := uuid.New()
	uuidstring := getHexUUID(epuuid)

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
	err = json.Unmarshal(jsonBlob, &ytd)
	if err != nil {
		log.Printf("Error reading metadata on URL %s\n", videoURL)
		return episode{}, err
	}

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

	_, err = deleteOldEpisodes()

	if err != nil {
		log.Println(err)
	}

	return fetchedEp, nil
}

func deleteOldEpisodes() (int64, error) {
	row := Database.QueryRow("SELECT sum(size) FROM episodes")
	var diskUsage int64
	err := row.Scan(&diskUsage)

	if err != nil {
		return -1, err
	}

	twoGigabytes := int64(2000000000)

	if diskUsage >= twoGigabytes {
		row := Database.QueryRow("SELECT uuid FROM episodes ORDER BY addedDate ASC LIMIT 1")
		var oldestUUID uuid.UUID
		err := row.Scan(&oldestUUID)

		if err != nil {
			return diskUsage, err
		}

		UUIDpath := filepath.Join(PodDirectory, getHexUUID(oldestUUID)) + ".mp3"
		err = os.Remove(UUIDpath)

		if err != nil {
			return diskUsage, err
		}

		_, err = Database.Exec("UPDATE episodes SET available=FALSE WHERE uuid=?", oldestUUID)

		if err != nil {
			return -1, err
		}

		return deleteOldEpisodes()
	}

	return diskUsage, nil
}

// This is basically a wrapper around youtube-dl, taking a URL
// and filename (without an extension) and returning the video
// metadata (or error), plus placing the audio in an mp3 file
// at the passed path
func getAndConvert(videoURL string, name string) ([]byte, int64, error) {
	command := exec.Command("youtube-dl",
		"-x", "--audio-format", "mp3",
		// E(-x)tract audio in mp3 format
		"--audio-quality", "64k",
		// 64kbps baby
		"--max-filesize", "1G",
		// Don't download files bigger than 1G
		"--print-json",
		// Print metadata to stdio
		// "--add-metadata", "--embed-thumbnail",
		// Embed metadata in mp3 and thumbnail
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
		log.Println("Error on " + name)
		log.Println(string(output))
		log.Println(err)
		return output, 0, err
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
