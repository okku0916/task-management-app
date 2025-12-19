package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, Task Management API!")
    })
    fmt.Println("Server starting at :8080...")
    http.ListenAndServe(":8080", nil)
}