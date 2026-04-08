package config

import (
	"html/template"
	"log"
	"path/filepath"
)

var Templates map[string]*template.Template

func TemplateParse() error {
	Templates = make(map[string]*template.Template)

	// Parse all .html files from views directory using glob
	files, err := filepath.Glob("./views/*.html")
	if err != nil {
		log.Printf("Error finding template files: %v\n", err)
		return err
	}

	// Parse each template file
	for _, file := range files {
		tmpl, err := template.ParseFiles(file)
		if err != nil {
			log.Printf("Error parsing template %s: %v\n", file, err)
			return err
		}
		// Use just the filename as the key
		name := filepath.Base(file)
		Templates[name] = tmpl
	}

	log.Printf("Parsed %d templates successfully\n", len(Templates))
	return nil
}

// GetTemplate retrieves a template by name
func GetTemplate(name string) *template.Template {
	if tmpl, exists := Templates[name]; exists {
		return tmpl
	}
	log.Printf("Template %s not found\n", name)
	return nil
}
