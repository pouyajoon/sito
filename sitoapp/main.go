package main

import (
	// "encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/websocket"
	"html/template"
	// "io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024 * 1024
)

type Page struct {
	Title string
	Body  []byte
}

func publicHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path[len("/public/"):]
	body, err := ioutil.ReadFile("public/" + filePath)
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
	log.Info("CLOSE CONNECTION")
	c.ws.WriteMessage(websocket.CloseMessage, []byte{})
	h.unregister <- c
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
			log.Info("read msg", msg)
			if msg.Player != "" {
				log.Info("read json", msg, err, msg.Player)
				// _, _ := json.Marshal(msg)
				h.broadcast <- "player" + msg.Player
				// c.ws.WriteJSON(msg)
			}
		}
		// mt, data, err := ws.ReadMessage()

		// ctx := log.Fields{"mt": mt, "data": data, "err": err}
		// if err != nil {
		// 	if err == io.EOF {
		// 		log.WithFields(ctx).Info("Websocket closed!")
		// 	} else {
		// 		log.Info(err)
		// 		log.WithFields(ctx).Error("Error reading websocket message")
		// 	}
		// 	break
		// }
		// switch mt {
		// case websocket.TextMessage:

		// 	err := json.Unmarshal(data, &msg)
		// 	if err != nil {
		// 		ctx["msg"] = msg
		// 		ctx["err"] = err
		// 		log.WithFields(ctx).Error("Invalid Message")
		// 		break
		// 	}
		// 	log.Info(msg, msg.Text, msg, msg.X, msg.Y)
		// default:
		// 	log.WithFields(ctx).Warning("Unknown Message!")
		// }
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
	}

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
	p := &Page{Title: "sito"}
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, p)
}

func main() {
	fmt.Println("sito")
	go h.run()

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
