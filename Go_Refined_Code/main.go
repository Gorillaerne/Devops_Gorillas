package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"devops_gorillas/database"
	"devops_gorillas/handlers"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "test endpoints")
}

func main() {
	// 1️⃣ Database
	if err := database.Connect(); err != nil {
		log.Fatal(err)
	}

	// 2️⃣ Templates
	tmpl := template.Must(
		template.ParseFiles(
			"templates/layout.html",
			"templates/search.html",
			"templates/about.html",
		),
	)

	// 3️⃣ Router
	r := mux.NewRouter()

	// Web
	r.HandleFunc("/", handlers.SearchHandler(tmpl)).Methods("GET")
	r.HandleFunc("/weather", homeHandler).Methods("GET")
	r.HandleFunc("/register", homeHandler).Methods("GET")
	r.HandleFunc("/login", homeHandler).Methods("GET")
	r.HandleFunc("/about", handlers.AboutHandler(tmpl)).Methods("GET")

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
