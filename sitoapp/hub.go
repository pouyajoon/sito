package main

import (
	// "encoding/json" V
	log "sito/sitoapp/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

type hub struct {
	// Registered clients
	clients map[*client]bool

	// Inbound messages
	broadcast chan string

	// Register requests
	register chan *client

	// Unregister requests
	unregister chan *client

	messages map[string]message

	content string
	id      int
}

var h = hub{
	broadcast:  make(chan string),
	register:   make(chan *client),
	unregister: make(chan *client),
	clients:    make(map[*client]bool),
	messages:   make(map[string]message),
	content:    "",
	id:         0,
}

func (h *hub) run() {
	log.Info("hub run")
	for {
		// log.Info("hub for")
		// h.broadcast <- string("salut")
		select {
		case c := <-h.register:
			log.Info("register")
			h.clients[c] = true
			c.send <- []byte(h.content)
			break
		case c := <-h.unregister:
			log.Info("unregister")
			_, ok := h.clients[c]
			if ok {
				delete(h.clients, c)
				close(c.send)
			}
			break
		case m := <-h.broadcast:
			// log.Info("broadcast case")
			h.content = m
			h.broadcastMessage()
			break
		default:
			// log.Info("no activity")
		}
	}
}

func (h *hub) broadcastMessage() {
	// log.Info("broadcast message", h.clients)
	for c := range h.clients {
		c.ws.WriteJSON(h.content)
		select {
		case c.send <- []byte(h.content):
			break

		// We can't reach the client
		default:
			close(c.send)
			delete(h.clients, c)
		}
	}
}
