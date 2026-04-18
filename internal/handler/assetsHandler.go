package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/static/")
	filePath := filepath.Clean("./internal/static/" + path)

	if _, err := os.Stat(filePath); err != nil {
		HandleError(w, "Not Found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, filePath)
}

func UploadsHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/uploads/")
	filePath := filepath.Clean("./uploads/" + path)

	if _, err := os.Stat(filePath); err != nil {
		HandleError(w, "Not Found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, filePath)
}
