package handlers

import (
	"html/template"
	"net/http"
)

func AboutHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "about.html", map[string]any{
			"content": "This is the about page.",
		})
	}
}
