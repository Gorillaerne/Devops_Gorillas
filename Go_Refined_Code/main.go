package main

import (
	"Refined_Project/handlers"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "test endpoints")
}


func main() {
    r := mux.NewRouter()

    r.HandleFunc("/", handlers.ServeLayout).Methods("GET")

	r.HandleFunc("/weather", homeHandler).Methods("GET")

	r.HandleFunc("/register", homeHandler).Methods("GET")

	r.HandleFunc("/login", homeHandler).Methods("GET")


	// API
	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/search", homeHandler).Methods("GET")

	api.HandleFunc("/weather", homeHandler).Methods("GET")

	api.HandleFunc("/register", homeHandler).Methods("POST")

	api.HandleFunc("/login", homeHandler).Methods("POST")

	api.HandleFunc("/logout", homeHandler).Methods("GET")


    http.ListenAndServe(":8080", r)
}