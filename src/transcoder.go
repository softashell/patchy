package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func computeHash(filePath string) ([]byte, error) {
	var result []byte

	file, err := os.Open(filePath)
	if err != nil {
		return result, err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return result, err
	}

	return hash.Sum(result), nil
}

func transcode(srcPath string, dstPath string) {
	cmd := "ffmpeg"

	dstPath = filepath.Join("static", dstPath)

	args := []string{"-i", srcPath, "-y", "-threads", "0", "-acodec", "libopus", "-b:a", "128k", "-vbr", "on", dstPath}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error, could not transcode song! Additional info: "+err.Error())
		os.Exit(1)
	}
}

func transcodeSong(filePath string) string {
	hash, err := computeHash(filePath)
	if err != nil {
		fmt.Printf("Unable to hash file: %v", err)
	}

	cacheFile := fmt.Sprintf("queue/%x.opus", hash)

	if !isCached(cacheFile) {
		transcode(filePath, cacheFile)
		// TODO: Remove old songs from cache
	} else {
		fmt.Println("Cache hit!")
	}

	return cacheFile
}

func isCached(filePath string) bool {
	filePath = filepath.Join("static", filePath)

	return exists(filePath)
}

func clearCache() {
	// TODO
}
