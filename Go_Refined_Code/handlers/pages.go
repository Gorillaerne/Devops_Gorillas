package handlers

import (
	"net/http"
)


func ServeLayout(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/layout.html")
}