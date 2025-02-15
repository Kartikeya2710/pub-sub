package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

// handlers and/or receives the operations that are to be performed on the actual subscriber over a websocket connection
type Subscriber struct {
	id string
	// The websocket connection to the subscriber
	conn *websocket.Conn
	// capacity before the subscriber is deemed too slow and absolutely demolished
	bufferSize int
	// The buffered channel to keep things that have been sent by  topic(s)
	ch chan string
}

func NewSubscriber(id string, conn *websocket.Conn, bufferSize int) *Subscriber {
	return &Subscriber{
		id:         id,
		conn:       conn,
		bufferSize: bufferSize,
		ch:         make(chan string, bufferSize),
	}
}

func (s *Subscriber) StartSubscriber() error {
	for {
		if len(s.ch) == cap(s.ch) {
			return fmt.Errorf("%s's channel is full. Annihilating it", s.id)
		}

		msg, ok := <-s.ch

		if !ok {
			return fmt.Errorf("%s's channel has been closed", s.id)
		}

		// send the message to the subscriber over the websocket connection
		s.conn.WriteMessage(websocket.TextMessage, []byte(msg))
	}
}
