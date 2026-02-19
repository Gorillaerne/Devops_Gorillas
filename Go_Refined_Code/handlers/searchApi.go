package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"devops_gorillas/database"
)

type SearchResult struct {
	Title   string `json:"title"`
	URL     string
	Content string `json:"content"`
}

type SearchResponse struct {
	Data []SearchResult `json:"data"`
}

func SearchAPIHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	log.Println("API search query:", q)

	var results []SearchResult

	if q != "" {
		rows, err := database.DB.Query(`
			SELECT title, content, url
			FROM pages
			WHERE title LIKE ? OR content LIKE ?
			LIMIT 20
		`, "%"+q+"%", "%"+q+"%")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var r SearchResult
			rows.Scan(&r.Title, &r.Content, &r.URL)
			results = append(results, r)
		}
	}

	response := SearchResponse{
		Data: results,
	}
	log.Println("WRAPPED RESPONSE EXECUTED")
	// 🟢 JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
