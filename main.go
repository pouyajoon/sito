package main

import "net/http"


func main() {
    println("sito")
    http.HandleFunc("/", someFunc)
    http.ListenAndServe(":8080", nil)
}

func someFunc(w http.ResponseWriter, req *http.Request) {
    w.Write([]byte("sito is cool"))
}
