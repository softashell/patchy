package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type request struct {
	File string
}

func handleSongs(utaChan chan string, reChan chan string, l *library, h *hub, q *queue) {
	ctChan := make(chan int)
	q.playing = false
	lastTime := 0
	stalled := false
	for msg := range utaChan {
		fmt.Println(msg)
		//Trigger this if there is a new song to be played
		if msg == "ns" && !q.playing {
			q.playing = true
			var ns *qsong

			//We'll want to make this concurrent since otherwise any requests which come in during the meanwhile will get pushed on wait
			//Needs some thought as to whether or not this is the best way to do stuff
			go func() {
				//Wait until all transcoding is done
				for q.transcoding {
				}

				if stalled {
					// First song was transcoded manually and player is ready to continue
					stalled = false
				}

				//Precondition: q.queue has at least 1 item in it.
				//Consume item in queue, if there's anything left, initiate a transcode
				ns = q.consume()
				lastTime = ns.Length

				if ns.Transcoded == "" {
					fmt.Println("Current song was not transcoded!")
				}

				msg := map[string]string{"cmd": "done", "File": ns.Transcoded}
				jsonMsg, _ := json.Marshal(msg)
				h.broadcast <- []byte(jsonMsg)
				fmt.Println("Sent done msg to clients")

				//Wait 4 seconds for clients to load the next song if necessary, then resume next song
				time.Sleep(4 * time.Second)
				fmt.Println("Sending NS to clients")

				//Tell clients to begin the song
				msg = map[string]string{"cmd": "NS", "Title": ns.Title, "Artist": ns.Artist, "Album": ns.Album, "Cover": "/art/" + GetAlbumDir(ns.File), "Time": strconv.Itoa(ns.Length)}
				jsonMsg, _ = json.Marshal(msg)
				h.broadcast <- []byte(jsonMsg)

				go timer(ns.Length, utaChan, ctChan)

				if len(q.queue) > 0 {
					fmt.Println("Queue has one or more items, performing next transcode in background")
					q.transcoding = true
					go q.transcodeNext()
				}
			}()
		}

		//If a song just finished, load in the next thing from queue if available
		if msg == "done" {
			q.playing = false
			songs := len(q.queue)

			if songs < 2 {
				if songs <= 0 {
					// Player was completely stopped and will not convert next song automatically
					stalled = true
				}

				// After current song ends it will try to keep at least one song playing and one more in queue
				added := q.addRandomSongs(2 - len(q.queue))

				// Broadcast added songs
				for _, song := range added {
					msg := map[string]string{"cmd": "queue", "Title": song.Title, "Artist": song.Artist}
					jsonMsg, _ := json.Marshal(msg)
					h.broadcast <- []byte(jsonMsg)
				}

				if len(q.queue) >= 1 && stalled {
					fmt.Println("Player stalled, transcoding next song manually")

					// Begin transcoding first song
					q.transcoding = true
					q.transcodeNext()
				}
			}

			// Play next song
			go func() {
				utaChan <- "ns"
			}()
		}

		if msg == "ctime" {
			if q.playing {
				ctChan <- 0
				reChan <- strconv.Itoa(<-ctChan)
			} else {
				//We want to actually do 100% here, do it later >.>
				reChan <- strconv.Itoa(lastTime)
			}
		}
		/*
			if msg == "queue" {
				jsonMsg, err := json.Marshal(q.queue)
				if err != nil {
					fmt.Println("Warning, could not jsonify queue")
				}
				utaChan <- string(jsonMsg)
			} else {
				utaChan <- ""
			}
		*/
	}
}

func handleRequests(requests chan *request, utaChan chan string, q *queue, l *library, h *hub) {
	for req := range requests {
		song, err := l.reqSearch(req.File)
		if err != nil {
			fmt.Println("Couldn't add request error: " + err.Error())
		} else {
			q.add(song)
			msg := map[string]string{"cmd": "queue", "Title": song.Title, "Artist": song.Artist}
			jsonMsg, _ := json.Marshal(msg)
			h.broadcast <- []byte(jsonMsg)

			if len(q.queue) == 1 {
				//This is safe to do because the loop guarentees that NS won't start until transcoding is finished
				fmt.Println("Queue has only one item, performing transcode")
				//Need a better way of doing this -- if we don't set q.transcoding to true prior to running the goroutine,
				//the songhandler can begin the NS functions before the transcodenext is able to set transcoding to true,
				//resulting in panics/errors
				q.transcoding = true
				go q.transcodeNext()
				if !q.playing {
					utaChan <- "ns"
				}
			}
		}
	}
}
