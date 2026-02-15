package main

import (
	"devops_gorillas/database"
	"devops_gorillas/handlers"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

  // Just for testing endpoints
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "test endpoints")
}

func main() {
	// 1️⃣ Database
	if err := database.Connect(); err != nil {
		log.Fatal(err)
	}


	// 3️⃣ Router
	r := mux.NewRouter()
	// API
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/search", handlers.SearchAPIHandler).Methods("GET")
	api.HandleFunc("/weather", homeHandler).Methods("GET")
	api.HandleFunc("/register", homeHandler).Methods("POST")
	api.HandleFunc("/login", homeHandler).Methods("POST")
	api.HandleFunc("/logout", homeHandler).Methods("GET")

	// 4️⃣ Server
	log.Println("Server kører på http://localhost:8080")
	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static"))),
	)
	log.Fatal(http.ListenAndServe(":8080", r))
}
