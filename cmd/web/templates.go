package main

import (
	"path/filepath"
	"text/template"
	"time"

	"lucasvinibox.isaacszf.net/internal/models"
)

// templateData is used to act as a holding structure for any dynamic data
// that we want to pass to out HTML templates
type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}

// Returns a nicely formatted string representation of a time.Time object
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// This is needed to use templa functions
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache
	cache := map[string]*template.Template{}

	// Using "filepath.Glob" to get all filepaths that match the pattern
	// ".ui/html/pages/*.tmpl.html"
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	// Loop through the page filepaths one-by-one
	for _, page := range pages {
		// Getting the file name (Ex: 'base.tmpl.html')
		name := filepath.Base(page)

		// The template.FuncMap must be registered with the template set before
		// you call the ParseFiles() method
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}

		// Calling "Parse.Glob" to add any components and pages
		ts, err = ts.ParseGlob("./ui/html/components/*.tmpl.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
