package handlers

import (
	"forum/config"
	"net/http"
)

func HandleError(w http.ResponseWriter, message string, status int) {
	w.WriteHeader(status)
	data := map[string]string{
		"Error":  message,
		"Status": http.StatusText(status),
	}
	config.RenderTemplate(w, "error.html", data)
}
