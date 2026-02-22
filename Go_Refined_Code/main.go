package main

import (
	"devops_gorillas/database"
	apiHandlers "devops_gorillas/handlers"
	"fmt"
	"github.com/gorilla/mux"
	cors "github.com/gorilla/handlers"
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

	http.ListenAndServe(":8080",
		cors.CORS(originsOk, headersOk, methodsOk)(r))
		log.Println("Server kører på http://localhost:8080")
}
