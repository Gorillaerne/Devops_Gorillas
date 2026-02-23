// Package main provides the entry point for the API server,
package main

import (
	"devops_gorillas/database"
	apiHandlers "devops_gorillas/handlers"
	"fmt"
	cors "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

// Just for testing endpoints
func homeHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(w, "test endpoints")
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
	api.HandleFunc("/search", apiHandlers.SearchAPIHandler).Methods("GET")
	api.HandleFunc("/weather", homeHandler).Methods("GET")
	api.HandleFunc("/register", apiHandlers.HandleApiRegister(database.DB)).Methods("POST")
	api.HandleFunc("/login", apiHandlers.HandleApiLogin(database.DB)).Methods("POST")
	api.HandleFunc("/logout", homeHandler).Methods("GET")

	// 4️⃣ Server
	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static"))),
	)
	// CORS options
	headersOk := cors.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := cors.AllowedOrigins([]string{"*"})
	methodsOk := cors.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})

	srv := &http.Server{
		Addr: ":8080",
		// The handler includes your router and CORS middleware
		Handler:      cors.CORS(originsOk, headersOk, methodsOk)(r),
		ReadTimeout:  5 * time.Second,   // Max time to read the request header/body
		WriteTimeout: 10 * time.Second,  // Max time to write the response
		IdleTimeout:  120 * time.Second, // Max time to keep an idle connection open
	}
	log.Println("Server kører på http://localhost:8080")
	// Use the srv instance to start the server
	
	log.Fatal(srv.ListenAndServe())
}
