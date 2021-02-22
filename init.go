package main

import (
	"fmt"
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
