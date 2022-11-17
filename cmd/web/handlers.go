package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// This is a handler (Handler = MVC Controller)
// This is also *application method, this is used so that we can use dependency injection
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Checking if req url is equal to "/"
	if r.URL.Path != "/" {
		// This function returns a 404 response (using custom helper)
		app.notFound(w)
		return
	}

	// Slice that contains the paths to the templates. The "base" template must be the
	// first one inside the slice
	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
		"./ui/html/components/nav.tmpl.html",
	}

	// Reading the files and storing the templates in a template set
	// (using "html/template" package)
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Using "ExecuteTemplate()" method to write the content inside "base" to the resp
	// body. The third parameter represents any dynamic data that we want to pass in,
	// which in this case is "nil" (by the time i'm writing this)
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, err)
	}
}

// Handler to view a snippet
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// localhost:4000/snippet/view?id=1

	// Extracting the ID value from URL (query parameter) and converting it
	// to a positive integer number (if user input is correct)
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		// Returning 404 if error
		app.notFound(w)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

// Handler to create a snippet
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	// Checking if request method is POST or not
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost) // Telling the user which methods are allowed

		// Returning status code 405 (method not allowed) if the request method
		// is different than POST
		app.clientError(w, 405)

		return
	}

	w.Write([]byte("Create a new snippet..."))
}
