package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"lucasvinibox.isaacszf.net/internal/models"
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

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
	}

	// Template data using method
	data := app.newTemplateData(r)
	data.Snippets = snippets

	// Using the helper method
	app.render(w, 200, "home.tmpl.html", data)
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

	// Retrieving the data for a specific record based on its ID
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}

		return
	}

	// Using method
	data := app.newTemplateData(r)
	data.Snippet = snippet

	// Using the helper method
	app.render(w, 200, "view.tmpl.html", data)
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

	// Dummy data to pass to the database
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
	}

	// Redirect the user to the relevant page for the snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), 303)
}
