package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
)

type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	// The hub.
	h *hub
}

//Reads in requests from the clients
func (c *connection) reader() {
	var msg string
	for {
		if err := websocket.Message.Receive(c.ws, &msg); err != nil {
			break
		}
		c.h.broadcast <- []byte(msg)
	}
	c.ws.Close()
}

//Sends broadcasts to clients
func (c *connection) writer() {
	for message := range c.send {
		fmt.Println("Sending a message to a connection")
		err := websocket.Message.Send(c.ws, string(message))
		if err != nil {
			break
		}
	}
	c.ws.Close()
}

//Socket handler -- Creates a new connection for each client
func handleSocket(ws *websocket.Conn, hub *hub) {
	c := &connection{send: make(chan []byte, 256), ws: ws, h: hub}
	c.h.register <- c
	defer func() { c.h.unregister <- c }()
	go c.writer()
	c.reader()
}
