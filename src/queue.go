package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type qsong struct {
	Title  string
	Album  string
	Artist string
	Length int
	File   string
}

type queue struct {
	// Library
	library *library
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
}

//Create a new queue
func newQueue(l *library) *queue {
	return &queue{
		library:     l,
		queue:       make([]*qsong, 0),
		CFile:       1,
		np:          nil,
		transcoding: false,
		pt:          "",
	}
}

//Consumes and returns the first value in the queue
//Precondition: queue has at least one item in it
func (q *queue) consume() *qsong {
	s := q.queue[0]

	if len(q.queue) < 1 {
		fmt.Println("Nothing in queue! Adding random song!")

		q.addRandom()
	}

	q.queue = q.queue[1:]

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
	q.queue = append(q.queue, s)
}

func (q *queue) addRandom() {
	fmt.Println("Adding random song")

	rand.Seed(time.Now().UnixNano())

	r := rand.Intn(len(q.library.library))
	s := q.library.library[r]

	fmt.Println("Chose song", r)

	length, err := strconv.Atoi(s["Length"])

	if err != nil {
		length = 0
	}

	q.add(&qsong{s["Title"], s["Album"], s["Artist"], length, s["file"]})
}

//Transcodes the next appropriate song
func (q *queue) transcodeNext() {
	//Need a better way of doing this -- perhaps transfer
	//From a nontranscoded queue to a transcoded queue?
	if q.pt == q.queue[0].File {
		fmt.Println("This song has already been transcoded!")
		return
	}
	q.pt = q.queue[0].File

	fmt.Println("Transcoding Song: ", q.queue[0].File)
	//We want to set this in whatever function calls transcodeNext because this
	//function is always called as a goroutine
	//q.transcoding = true
	transcode(musicDir + "/" + q.queue[0].File)
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
