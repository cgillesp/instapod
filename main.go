package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sys/unix"
)

var podDirectory string

type episode struct {
	title       string
	description string
	uuid        uuid.UUID
	addDate     time.Time
	pubDate     time.Time
}

func main() {
	checkDeps()
	podDir := getDir()
	uuidbin, _ := uuid.New().MarshalBinary()
	uuidstring := hex.EncodeToString(uuidbin)

	fullpath := filepath.Join(podDir, uuidstring)

	if len(os.Args) < 2 {
		fmt.Println("Please supply a URL as an argument")
		os.Exit(1)
	}

	url := os.Args[1]
	getAndConvert(url, fullpath)
}

func getDir() string {
	usr, _ := user.Current()
	homedir := usr.HomeDir
	podDir := filepath.Join(homedir, ".instapod")

	_, err := os.Stat(podDir)

	if os.IsNotExist(err) {
		err := os.Mkdir(podDir, 0755)
		if err != nil {
			fmt.Println("Could not create ~/.instapod folder. Check your permissions.")
		}
	} else if err != nil {
		panic(err)
	}

	if unix.Access(podDir, unix.W_OK) != nil {
		fmt.Println("Could not read ~/.instapod folder. Check your permissions.")
		os.Exit(1)
	}

	return podDir
}

func getAndConvert(videoURL string, name string) error {
	command := exec.Command("youtube-dl", "-x", "--audio-format", "mp3",
		"--audio-quality", "64k", "--embed-thumbnail",
		"--add-metadata", "--print-json",
		"--postprocessor-args", "-ac 1",
		"-o", name+".%(ext)s",
		videoURL)
	output, err := command.Output()

	fmt.Print(string(output))

	if err != nil {
		panic(err)
	}

	return nil
}

func checkDeps() {
	// fmt.Print("Checking dependencies... ")
	// defer fmt.Print("done \n")
	_, err := exec.Command("youtube-dl", "--version").Output()

	if err != nil {
		fmt.Println(
			"Instapod requires youtube-dl (https://github.com/ytdl-org/youtube-dl)")
		os.Exit(1)
	}

	_, err = exec.Command("ffmpeg", "-version").Output()

	if err != nil {
		fmt.Println("Instapod requires FFmpeg (https://ffmpeg.org)")
		os.Exit(1)
	}
}
