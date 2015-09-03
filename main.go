package main

import (
    "fmt"
    "net/http"
    "os"
)

func main() {
    fmt.Println("sito")
    http.HandleFunc("/", someFunc)
    http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}

func someFunc(res http.ResponseWriter, req *http.Request) {
    fmt.Fprintln(res, "sito is cool")

}
