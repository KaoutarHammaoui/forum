package handlers

import (
	"forum/config"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	config.RenderTemplate(w, "home.html", nil)
}