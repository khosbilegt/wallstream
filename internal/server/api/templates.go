package api

import (
	"html/template"
	"log"
	"path/filepath"
)

var templates *template.Template

func LoadTemplates() {
	var err error
	templates, err = template.ParseGlob(filepath.Join("internal/server/templates", "*.html"))
	log.Println("Loading templates from", filepath.Join("internal/server/templates", "*.html"))
	if err != nil {
		log.Fatalf("failed to load templates: %v", err)
	}
}
