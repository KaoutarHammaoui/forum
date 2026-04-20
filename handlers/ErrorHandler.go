package handlers

import (
	"forum/config"
	"net/http"
)

func HandleError(w http.ResponseWriter, message string, status int) {
	w.WriteHeader(status)
	data := homeError{
		Error:  message,
		Status: status,
	}
	config.RenderTemplate(w, "error.html", data)
}
