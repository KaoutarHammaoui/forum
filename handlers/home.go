package handlers

import (
	"forum/config"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	tmpl := config.GetTemplate("home.html")
	if tmpl == nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}

