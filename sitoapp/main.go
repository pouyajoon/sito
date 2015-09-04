package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"html/template"
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

func handleWebsocket(w http.ResponseWriter, r *http.Request) {}

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
