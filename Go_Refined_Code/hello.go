package main

import (
    "fmt"
    "net/http"
    "github.com/gorilla/mux" 
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Gorilla is running!")
}

func main() {
    r := mux.NewRouter()

    r.HandleFunc("/", homeHandler)

    fmt.Println("Server starting on :8080")
    http.ListenAndServe(":8080", r) 
}