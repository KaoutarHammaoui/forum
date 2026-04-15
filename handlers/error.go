package handlers

import (
	"html/template"
	"net/http"
)

func Error(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)

	tmpl, err := template.ParseFiles("./web/templates/error.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := map[string]string{
		"Error":   message,
		"Status": http.StatusText(status),
	}

	tmpl.Execute(w, data)
}
