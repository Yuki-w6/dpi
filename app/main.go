package main

import (
    "net/http"
    "log"
)

func main() {
    fs := http.FileServer(http.Dir("./build"))
    http.Handle("/", fs)

    log.Println("Serving on http://localhost:1000")
    err := http.ListenAndServe(":1000", nil)
    if err != nil {
        log.Fatal(err)
    }
}