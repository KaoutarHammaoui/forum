package handlers

import (
	"net/http"
	"os"
	"strings"
)

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/static" || r.URL.Path == "/static/" {
		HandleError(w, "Forbidden", http.StatusForbidden)
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/static/")
	filePath := "./static/" + path
	file, err := os.Stat(filePath)
	if err != nil {
		HandleError(w, "Not Found", http.StatusNotFound)
		return
	}
	if file.IsDir() {
		HandleError(w, "Forbidden", http.StatusForbidden)
		return
	}

	http.ServeFile(w, r, filePath)
}

func UploadsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/uploads" || r.URL.Path == "/uploads/" {
		HandleError(w, "Forbidden", http.StatusForbidden)
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/uploads/")
	filePath := "./uploads/" + path
	file, err := os.Stat(filePath)
	if err != nil {
		HandleError(w, "Not Found", http.StatusNotFound)
		return
	}
	if file.IsDir() {
		HandleError(w, "Forbidden", http.StatusForbidden)
		return
	}

	http.ServeFile(w, r, filePath)
}
