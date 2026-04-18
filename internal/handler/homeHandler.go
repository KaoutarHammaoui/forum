package handler

import (
	"forum/internal/config"
	"forum/internal/middleware"
	models "forum/internal/model"

	"net/http"
)

type Data struct {
	Title      string
	IsLoggedIn bool
	Categories []models.Category
}

func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		HandleError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	_, loggedIn := middleware.IsAuthenticated(r)

	data := Data{
		Title:      "Home",
		IsLoggedIn: loggedIn,
		Categories: []models.Category{}, // replace with real data later
	}

	// Example response (you'll likely use templates instead)
	config.RenderTemplate(w, "index.html", data)
}
