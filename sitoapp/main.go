package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/websocket"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// type Page struct {
// 	Title string
// 	Body  []byte
// }

// func main() {
// 	fmt.Println("sito")
// 	http.HandleFunc("/", indexHandler)
// 	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
// 	if err != nil {
// 		panic(err)
// 	}
// }

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func publicHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path[len("/public/"):]
	body, err := ioutil.ReadFile("sitoapp/public/" + filePath)
	if err == nil {
		fmt.Fprintf(w, string(body))
	}
}

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
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

	// id := rr.register(ws)

	for {
		mt, data, err := ws.ReadMessage()
		ctx := log.Fields{"mt": mt, "data": data, "err": err}
		if err != nil {
			if err == io.EOF {
				log.WithFields(ctx).Info("Websocket closed!")
			} else {
				log.WithFields(ctx).Error("Error reading websocket message")
			}
			break
		}
		switch mt {
		case websocket.TextMessage:
			msg, err := validateMessage(data)
			if err != nil {
				ctx["msg"] = msg
				ctx["err"] = err
				log.WithFields(ctx).Error("Invalid Message")
				break
			}
			rw.publish(data)
		default:
			log.WithFields(ctx).Warning("Unknown Message!")
		}
	}

	// rr.deRegister(id)

	ws.WriteMessage(websocket.CloseMessage, []byte{})
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	// title := r.URL.Path[len("/view/"):]
	p := &Page{Title: "sito"}
	t, _ := template.ParseFiles("sitoapp/templates/index.html")
	t.Execute(w, p)
}

func main() {
	fmt.Println("sito")
	// http.Handle("/public/", http.FileServer(http.Dir("./siteoapp/public")))

	port := os.Getenv("PORT")
	if port == "" {
		log.WithField("PORT", port).Fatal("$PORT must be set")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", handleWebsocket)
	mux.HandleFunc("/public/", publicHandler)
	mux.HandleFunc("/", viewHandler)

	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":" + port)

}
