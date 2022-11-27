package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

// This method helps write an error message and stack-trace to the errorLog
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, "Internal Server Error", 500)
}

// This method sends a specific status code and description to the user
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// Wrapper to 404 not found
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, 404)
}

// Helper method to reduce code duplication when rendering a template
func (app *application) render(
	w http.ResponseWriter,
	status int,
	page string,
	data *templateData) {
	// Getting the appropriate template set from the cache-map
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	// Initializing Buffer
	buf := new(bytes.Buffer)

	// Writing the template to the buffer, instead of straight to the http.ResponseWriter
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// If the template is written to the buffer without any errors, we are safe to go ahead
	w.WriteHeader(status)

	buf.WriteTo(w)
}

// Returns a pointer to a templateData struct initialized with the current year
func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}
