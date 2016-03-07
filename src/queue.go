package main

import (
	"fmt"
	"path/filepath"
	"strconv"
)

type qsong struct {
	Title      string
	Album      string
	Artist     string
	Length     int
	File       string
	Transcoded string
}

type queue struct {
	// Queue.
	queue []*qsong
	//Current playing song
	np *qsong
	//Transcoding status
	transcoding bool
	//Playing status
	playing bool
	//Library
	l *library
}

//Create a new queue
func newQueue(library *library) *queue {
	return &queue{
		queue:       make([]*qsong, 0),
		np:          nil,
		transcoding: false,
		l:           library,
	}
}

//Consumes and returns the first value in the queue
//Precondition: queue has at least one item in it
func (q *queue) consume() *qsong {
	if len(q.queue) < 1 {
		fmt.Println("Nothing left in queue!")
	}

	s := q.queue[0]

	if len(q.queue) > 1 {
		q.queue = q.queue[1:]
	} else {
		//There has to be a better way of doing this
		q.queue = make([]*qsong, 0)
	}

	q.np = s

	return s
}

//Adds a new item to the queue
func (q *queue) add(s *qsong) {
	fmt.Println("Added:", s.File)

	q.queue = append(q.queue, s)
}

// Gets random song from library, and adds it to queue
func (q *queue) addRandomSongs(count int) []qsong {
	fmt.Println("Adding", count, "random songs!")

	songs := []qsong{}

	for i := 0; i < count; i++ {
		song := q.l.getRandomSong()

		// TODO: Check for dupes

		st, err := strconv.Atoi(song["Time"])
		if err != nil {
			fmt.Println("Could not convert time to int for", song["file"])
			continue
		}

		songs = append(songs, qsong{song["Title"], song["Album"], song["Artist"], st, song["file"], ""})
	}

	for i, _ := range songs {
		q.add(&songs[i])
	}

	return songs
}

//Transcodes the next appropriate song
func (q *queue) transcodeNext() {
	s := q.queue[0]

	//Need a better way of doing this -- perhaps transfer
	//From a nontranscoded queue to a transcoded queue?
	if s.Transcoded != "" {
		fmt.Println("This song has already been transcoded!", s.File)
		return
	}

	fmt.Println("Transcoding song:", s.File, " > ")

	//We want to set this in whatever function calls transcodeNext because this
	//function is always called as a goroutine
	//q.transcoding = true
	songPath := filepath.Join(musicDir, s.File)
	s.Transcoded = transcodeSong(songPath)

	fmt.Println("transcode ok", s.Transcoded)

	q.transcoding = false
}
