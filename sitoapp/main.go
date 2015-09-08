package main

import (
	"encoding/json"
	"fmt"
	log "sito/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"sito/Godeps/_workspace/src/github.com/codegangsta/negroni"
	"sito/Godeps/_workspace/src/github.com/gorilla/websocket"
	// "html/template"
	// "io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024 * 1024
	// rootPath       = "sitoapp/"
	rootPath = ""
)

// var int id
// id = 0

type Page struct {
	Title string
	Body  []byte
}

func publicHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path[len("/public/"):]
	body, err := ioutil.ReadFile(rootPath + "public/" + filePath)
	if err == nil {
		fmt.Fprintf(w, string(body))
	}
}

// connection is an middleman between the websocket connection and the hub.
type client struct {
	// The websocket connection.
	ws *websocket.Conn
	// Buffered channel of outbound messages.
	send chan []byte
	id   int
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// message sent to us by the javascript client
type message struct {
	Player string `json:"player"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
}

func closeConnection(c *client) {
	log.Info("CLOSE CONNECTION", c.id)
	c.ws.WriteMessage(websocket.CloseMessage, []byte{})
	h.unregister <- c
	delete(h.messages, strconv.Itoa(c.id))
}

func handleMessage(c *client) {
	for {
		var msg message
		err := c.ws.ReadJSON(&msg)
		if err != nil {
			log.Error(err)
			closeConnection(c)
			return
		} else {
			if msg.Player != "" {
				h.messages[strconv.Itoa(c.id)] = msg
			}
		}
	}
	closeConnection(c)
}

func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.WithField("err", err).Println("Upgrading to websockets")
		http.Error(w, "Error Upgrading to websockets", 400)
		return
	}

	log.Info("NEW CONNECTION")

	c := &client{
		send: make(chan []byte, maxMessageSize),
		ws:   ws,
		id:   h.id,
	}
	h.id += 1
	h.register <- c
	// select {
	// case h.register <- c:
	// 	fmt.Println("sent message")
	// default:
	// 	fmt.Println("no message sent")
	// }

	go handleMessage(c)

}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadFile(rootPath + "templates/index.html")
	fmt.Fprintf(w, string(body))
	// p := &Page{Title: "sito"}
	// t, err := template.ParseFiles("templates/index.html")
	// if err != nil {
	// 	log.Error(err)
	// }
	// t.Execute(w, p)
}

func interval() {
	ticker := time.NewTicker(time.Millisecond * 50)
	go func() {
		for range ticker.C {
			s, _ := json.Marshal(h.messages)
			h.broadcast <- s
		}
	}()
}

func main() {
	fmt.Println("sito")
	go h.run()
	go interval()

	port := os.Getenv("PORT")
	if port == "" {
		log.WithField("PORT", port).Fatal("$PORT must be set")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", handleWebsocket)
	// http.Handle("/public", http.FileServer(http.Dir("/public")))
	mux.HandleFunc("/public/", publicHandler)
	mux.HandleFunc("/", viewHandler)

	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":" + port)

}
