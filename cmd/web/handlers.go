package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"lucasvinibox.isaacszf.net/internal/models"
)

// This is a handler (Handler = MVC Controller)
// This is also *application method, this is used so that we can use dependency injection
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Because httprouter matches the "/" path exactly, we can now remove the
	// manual check of "r.Url != '/'" from this handler

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
	// Collecting the parameters passed by the user using httprouter req context
	params := httprouter.ParamsFromContext(r.Context())

	// The method "ByName" return the value of "id" from the slice and
	// validate it as normal
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
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

// Handler to create snippet
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// Dummy data to pass to the database
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
	}

	// Redirect the user to the relevant page for the snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), 303)
}

// Handler to see the form that will be use to create a snippet
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {

}
