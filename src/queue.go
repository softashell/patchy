package main

import (
	"fmt"
	"os"
	"strconv"
)

type qsong struct {
	Title  string
	Album  string
	Artist string
	Length int
	File   string
}

type queue struct {
	// Queue.
	queue []*qsong
	//Current playing song
	np *qsong
	//Current file in use
	CFile int
	//Transcoding status
	transcoding bool
	//Playing status
	playing bool
	//Previously transcoded file -- Used to prevent dupes
	pt string
	//Library
	l *library
}

//Create a new queue
func newQueue(library *library) *queue {
	return &queue{
		queue:       make([]*qsong, 0),
		CFile:       1,
		np:          nil,
		transcoding: false,
		pt:          "",
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

	if q.CFile == 1 {
		q.CFile = 2
	} else {
		q.CFile = 1
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

		songs = append(songs, qsong{song["Title"], song["Album"], song["Artist"], st, song["file"]})
	}

	for _, song := range songs {
		q.add(&song)
	}

	return songs
}

//Transcodes the next appropriate song
func (q *queue) transcodeNext() {
	song := q.queue[0]

	//Need a better way of doing this -- perhaps transfer
	//From a nontranscoded queue to a transcoded queue?
	if q.pt == song.File {
		fmt.Println("This song has already been transcoded!", song.File)
		return
	}

	q.pt = song.File

	fmt.Println("Transcoding Song:", song.File)

	//We want to set this in whatever function calls transcodeNext because this
	//function is always called as a goroutine
	//q.transcoding = true
	transcode(musicDir + "/" + song.File)

	//Rename to opposite of current file, since the clients will be told to go
	//to the next song after this
	//Transcodes will happen BEFORE consumes(need to create the file for client use)
	//if there is only one thing in the queue
	//other wise transcodes occur afterwards(since you want to transcodein background)
	if q.CFile == 1 {
		fmt.Println("Renaming Song to ns2.opus")
		os.Rename("static/queue/next.opus", "static/queue/ns2.opus")
	} else {
		fmt.Println("Renaming Song to ns1.opus")
		os.Rename("static/queue/next.opus", "static/queue/ns1.opus")
	}

	q.transcoding = false
}
