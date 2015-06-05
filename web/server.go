package main

import (
    "log"
    "fmt"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/", fs)
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    log.Println("Listening...")
    http.ListenAndServe(":3000", nil)
}
