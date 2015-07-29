package main

import "os"

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
	np    *qsong
	CFile int
}

//Create a new queue
func newQueue() *queue {
	return &queue{
		queue: make([]*qsong, 0),
		CFile: 1,
		np:    nil,
	}
}

//Consumes and returns the first value in the queue
//Precondition: queue has at least one item in it
func (q *queue) consume() *qsong {
	/*
		if len(q.Queue) < 1 {
			return nil, errors.New("Nothing in queue!")
		}
	*/
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
	q.queue = append(q.queue, s)
}

//Transcodes the next appropriate song
func (q *queue) transcodeNext() {
	l := len(q.queue)
	switch {
	case l == 1:
		transcode(musicDir + "/" + q.queue[0].File)
	case l > 1:
		transcode(musicDir + "/" + q.queue[1].File)
	}
	//Rename to opposite of current file, since the clients will be told to go
	//to the next song after this
	//Transcodes will happen BEFORE consumes(need to create the file for client use)
	//if there is only one thing in the queue
	//other wise transcodes occur afterwards(since you want to transcodein background)
	if q.CFile == 1 {
		os.Rename("static/queue/next.mp3", "static/queue/ns2.mp3")
	} else {
		os.Rename("static/queue/next.mp3", "static/queue/ns1.mp3")
	}
}
