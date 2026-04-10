package config

import "net/http"

func LoadAssets() {
	// Serve static files (CSS, JS, images, etc.)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
}
