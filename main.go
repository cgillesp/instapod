package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {

	fmt.Print("Checking dependencies... ")
	defer fmt.Print("done \n")
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

	outp, err := exec.Command("youtube-dl", "-x", "--audio-format mp3",
		"--audio-quality 64k", "https://www.youtube.com/watch?v=SQ8UDzUvMg0").Output()
	if err != nil {
		panic(err)
	}
	fmt.Print(string(outp))
}
