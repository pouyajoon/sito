package main

import (
    "fmt"
    "net/http"
    "os"
)

func main() {
    fmt.Println("sito")
    http.HandleFunc("/", someFunc)
    err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
    if err != nil {
      panic(err)
    }
}

func someFunc(res http.ResponseWriter, req *http.Request) {
    fmt.Fprintln(res, "sito is cool")

}
