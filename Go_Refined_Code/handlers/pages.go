package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
	"log"
)

type Page struct {
    Title    string
    Content  string
    URL string
}

type SearchData struct {
    Q         string
    SearchResults []Page
}


func ServeLandingPage(db *sql.DB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")

		language := r.URL.Query().Get("language")
		if language == "" {
			language = "en" // Default to English
		}

		var results []Page
log.Println("Q: " + q)

		if q != "" {

		
            // Use placeholders (?) to prevent SQL Injection
            rows, err := db.Query("SELECT title, content, url FROM pages WHERE language = ? AND content LIKE ?", 
                language, "%"+q+"%")
            if err == nil {
                defer rows.Close()
                for rows.Next() {
                    var p Page
                    if err := rows.Scan(&p.Title, &p.Content, &p.URL); err == nil {
                        results = append(results, p)
                    }
                }
            }
        }

        // 3. Pass the data to the template
        data := SearchData{
            Q:         q,
            SearchResults: results,
        }

		tmpl.ExecuteTemplate(w, "layout.html", data)
	}
}

func ServeAboutPage(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "about.html", nil)
	}
}

func ServeLoginPage(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "layout.html", nil)
	}
}

func ServeRegisterPage(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "layout.html", nil)
	}
}
