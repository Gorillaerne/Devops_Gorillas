package handlers

import (
	"devops_gorillas/database"
	"html/template"
	"net/http"
)

type SearchResult struct {
	Title   string
	URL     string
	Content string
}

func SearchHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")

		var results []SearchResult

		if q != "" {
			rows, err := database.DB.Query(`
			SELECT title, content, url
			FROM pages
			WHERE title LIKE ? OR content LIKE ?
			LIMIT 20
`, "%"+q+"%", "%"+q+"%")

			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			defer rows.Close()

			for rows.Next() {
				var r SearchResult
				rows.Scan(&r.Title, &r.Content, &r.URL)
				results = append(results, r)
			}
		}

		tmpl.ExecuteTemplate(w, "layout.html", map[string]any{
			"search_results": results,
			"q":              q,
		})
	}
}
