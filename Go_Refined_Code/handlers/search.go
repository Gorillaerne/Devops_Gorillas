package

import "net/http"
func searchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")

	var results []SearchResult

	if q != "" {
		rows, err := database.DB.Query(`
			SELECT title, description, url
			FROM pages
			WHERE title LIKE ? OR description LIKE ?
			LIMIT 20
		`, "%"+q+"%", "%"+q+"%")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var r SearchResult
			rows.Scan(&r.Title, &r.Description, &r.URL)
			results = append(results, r)
		}
	}

	tmpl.ExecuteTemplate(w, "search.html", map[string]any{
		"search_results": results,
		"request": map[string]any{
			"args": map[string]string{
				"q": q,
			},
		},
	})
}
