package main

import (
	"fmt"
	"os"
	"os/exec"
)

func transcode(song string) {
	cmd := "ffmpeg"
	args := []string{"-i", song, "-y", "-threads", "0", "-acodec", "libopus", "-b:a", "128k", "-vbr", "on", "static/queue/next.opus"}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error, could not transcode song! Additional info: "+err.Error())
		os.Exit(1)
	}
}
