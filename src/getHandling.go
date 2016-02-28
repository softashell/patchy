package main

import (
	"encoding/json"
	"github.com/fhs/gompd/mpd"
	"github.com/hoisie/web"
	"net/http"
	"path/filepath"
	"strconv"
)

func getCover(ctx *web.Context, album string) string {
	dir := musicDir + "/" + album

	cover := "static/image/missing.png"

	cover_files := [...]string{
		"cover.jpg",
		"cover.png",
		"AlbumArt.jpg",
		"AlbumArt.png",
		"folder.jpg",
		"folder.png",
	}

	// Try the most generic album art names and return first existing match
	for _, filename := range cover_files {
		file := filepath.Join(dir, filename)

		if exists(file) {
			cover = file

			break
		}
	}

	http.ServeFile(ctx.ResponseWriter, ctx.Request, cover)

	return ""
}

func getSong(ctx *web.Context, song string) string {
	http.ServeFile(ctx.ResponseWriter, ctx.Request, "static/queue/"+song)

	return ""
}

func getSearchRes(ctx *web.Context, req string, l *library) string {
	res := l.asyncSearch(req)
	jsonMsg, _ := json.Marshal(res)

	return string(jsonMsg)
}

func getNowPlaying(ctx *web.Context, utaChan chan string, reChan chan string, queue *queue, listeners int) string {
	song := make(map[string]string)

	if np := queue.np; np != nil {
		utaChan <- "ctime"
		ctime := <-reChan

		utaChan <- "cfile"
		cfile := <-reChan

		song["Title"] = np.Title
		song["Artist"] = np.Artist
		song["Album"] = np.Album
		song["file"] = np.File
		song["Cover"] = GetAlbumDir(np.File)
		song["Time"] = strconv.Itoa(np.Length)

		song["ctime"] = ctime
		song["cfile"] = cfile
	} else {
		song["Title"] = "N/A"
		song["Artist"] = "N/A"
		song["Album"] = "N/A"
		song["file"] = "lol"
		song["Time"] = "0"

		song["ctime"] = "0"
		song["cfile"] = "1"
	}

	song["listeners"] = strconv.Itoa(listeners)

	jsonMsg, _ := json.Marshal(song)

	return string(jsonMsg)
}

func getLibrary(ctx *web.Context, subset []mpd.Attrs) string {
	jsonMsg, _ := json.Marshal(subset)

	return string(jsonMsg)
}

func getQueue(ctx *web.Context, q *queue, h *hub, utaChan chan string) string {
	//Let the song handler return a JSONify'd queue
	jsonMsg, _ := json.Marshal(q.queue)

	return string(jsonMsg)
}
