package config

import "net/http"

func LoadAssets() {
	// Serve static files (CSS, JS, images, etc.)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	// Sert les fichiers du dossier "uploads/" via l'URL "/uploads/"
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

}
