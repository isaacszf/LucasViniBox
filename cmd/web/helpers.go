package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
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
