package handler

import (
	"net/http"
	"os"
	"strings"
)

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	referer := r.Header.Get("Referer")
	if referer == "" {
		HandleError(w, "Forbidden", http.StatusForbidden)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/static/")
	filePath := "./static/" + path

	if _, err := os.Stat(filePath); err != nil {
		HandleError(w, "Not Found", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, filePath)
}

func UploadsHandler(w http.ResponseWriter, r *http.Request) {

	referer := r.Header.Get("Referer")
	if referer == "" {
		HandleError(w, "Forbidden", http.StatusForbidden)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/uploads/")
	filePath := "./uploads/" + path

	if _, err := os.Stat(filePath); err != nil {
		HandleError(w, "Not Found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, filePath)
}