package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"

	"golang.org/x/sys/unix"
)

// Determines whether the system has the neccesary dependencies
// to run the program
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

// Checks if the data directory exists and is writeable,
// creating it if it isn't and complaining otherwise
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

type configuration struct {
	AddKey       string
	addKeyBytes  []byte
	ReadKey      string
	readKeyBytes []byte
	Title        string
	Link         string
	Description  string
	BaseURL      string
	ImageURL     string
}

func getConfig() configuration {
	configPath := filepath.Join(PodDirectory, "config.json")
	configBytes, err := ioutil.ReadFile(configPath)

	if err != nil {
		if os.IsNotExist(err) {
			newConfigFile, err := json.Marshal(newConfig())
			ioutil.WriteFile(configPath, newConfigFile, 0744)
			if err != nil {
				panic(err)
			}
			return (getConfig())
		}

		panic(err)
	}

	config := configuration{}

	json.Unmarshal(configBytes, &config)

	fmt.Printf("Add Key: %s\n", string(config.AddKey))
	fmt.Printf("Read Key: %s\n", string(config.ReadKey))

	// we need these for constant time comparisons
	config.addKeyBytes = []byte(config.AddKey)
	config.readKeyBytes = []byte(config.ReadKey)

	return config
}

func newConfig() configuration {
	return configuration{
		AddKey:      getKeyBase64(),
		ReadKey:     getKeyBase64(),
		Title:       "Instacast feed",
		Description: "Catch up on whatever you've been meaning to listen to!",
		Link:        "demos.charliegillespie.com",
		BaseURL:     "https://demos.charliegillespie.com/",
		ImageURL:    "",
	}
}

func getKeyBase64() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b)
}
