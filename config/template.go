package config

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var Templates map[string]*template.Template

func TemplateParse() error {
	Templates = make(map[string]*template.Template)
	pages, err := filepath.Glob("./views/*.html")
	if err != nil {
		return err
	}
	for _, page := range pages {
		// Ne pas parser partials.html comme page principale
		if filepath.Base(page) == "composants.html" {
			continue
		}

		// Parser la page AVEC le fichier partials
		tmpl, err := template.ParseFiles(page, "./views/composants.html")
		if err != nil {
			log.Printf("Error parsing template %s: %v\n", page, err)
			return err
		}

		name := filepath.Base(page)
		Templates[name] = tmpl
	}

	log.Printf("Parsed %d templates successfully\n", len(Templates))
	return nil
}

func GetTemplate(name string) *template.Template {
	if tmpl, exists := Templates[name]; exists {
		return tmpl
	}
	log.Printf("Template %s not found\n", name)
	return nil
}

func RenderTemplate(w http.ResponseWriter, name string, data any) {
	tmpl := GetTemplate(name)
	if tmpl == nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Printf("Error executing template %s: %v\n", name, err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
