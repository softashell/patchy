package main

import (
	"errors"
	"fmt"
	"github.com/fhs/gompd/mpd"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type library struct {
	// The library.
	library []mpd.Attrs
}

//Create a new queue
func newLibrary() *library {
	fmt.Println("Connecting to MPD")

	conn, err := mpd.Dial("tcp", "127.0.0.1:6600")
	if err != nil {
		fmt.Println("Error: could not connect to MPD, exiting.", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	songs, err := conn.ListAllInfo("/")
	if err != nil {
		fmt.Println("Error: could not list MPD library, exiting.", err.Error())
		os.Exit(1)
	}

	shuffle(songs)

	return &library{
		library: songs,
	}
}

//returns a small selection of songs for initial display
func (l *library) selection() []mpd.Attrs {
	songs := l.library[:15]

	for i, song := range songs {
		songs[i]["Cover"] = GetAlbumDir(song["file"])
	}

	return songs
}

//Searches for a request and returns the first song which matches
func (l *library) reqSearch(songPath string) (*qsong, error) {
	for _, song := range l.library {
		if song["file"] == songPath {
			//if song["Title"] == title && (song["Album"] == album || song["Artist"] == artist) {
			fmt.Println("Found song: " + song["file"])
			st, err := strconv.Atoi(song["Time"])
			if err != nil {
				fmt.Println("Couldn't get song due to time conversion error!")
				return nil, errors.New("Couldn't convert Time to int!")
			}
			return &qsong{Title: song["Title"], Album: song["Album"], Artist: song["Artist"], Length: st, File: song["file"]}, nil
		}
	}
	return nil, errors.New("No songs found!")
}

func (l *library) asyncSearch(req string) []mpd.Attrs {
	res := make([]mpd.Attrs, 0)

	req = strings.ToLower(req)

	//There has to be a faster way to do this >.>
	for _, song := range l.library {
		artist := strings.ToLower(song["Artist"])
		title := strings.ToLower(song["Title"])
		album := strings.ToLower(song["Album"])
		file := strings.ToLower(song["file"])

		if strings.Contains(title, req) || strings.Contains(album, req) || strings.Contains(artist, req) || strings.Contains(file, req) {
			song["Cover"] = GetAlbumDir(song["file"])
			res = append(res, song)
			if len(res) == 100 {
				break
			}
		}
	}

	return res
}

//Updates the library
func (l *library) update() error {
	var conn *mpd.Client

	fmt.Println("Connecting to MPD")
	conn, err := mpd.Dial("tcp", "127.0.0.1:6600")
	if err != nil {
		fmt.Println("Error: could not connect to MPD for lib update")
		return errors.New("Could not connect to MPD!")
	}
	defer conn.Close()

	_, err = conn.Update("")
	if err != nil {
		fmt.Println("Error: could not update library!")
		return err
	}

	//Let the update happen
	time.Sleep(2 * time.Second)
	songs, err := conn.ListAllInfo("/")
	if err != nil {
		fmt.Println("Error: could not retrieve new library!")
		return err
	}

	l.library = songs
	return nil
}

func (l *library) getRandomSong() mpd.Attrs {
	rand.Seed(time.Now().UnixNano())

	total := len(l.library)

	r := rand.Intn(total)
	s := l.library[r]

	return s
}
