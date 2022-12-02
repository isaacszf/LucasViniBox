package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/julienschmidt/httprouter"
	"lucasvinibox.isaacszf.net/internal/models"
)

// This represents the form data and validation errors for the form fields
type snippetCreateForm struct {
	Title       string
	Content     string
	Expires     int
	FieldErrors map[string]string
}

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
	// Calling r.ParseForm() which adds any data in POST req bodies to the r.PostForm
	// map
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, 400)
		return
	}

	// Parsing string to numbers. This is needed because "expires" needs to be a number
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, 400)
		return
	}

	// Creating an instance of the snippetCreateForm struct containing the values
	// from the form and an empty map for any validation errors
	form := snippetCreateForm{
		Title:       r.PostForm.Get("title"),
		Content:     r.PostForm.Get("content"),
		Expires:     expires,
		FieldErrors: map[string]string{},
	}

	// Verifying if the title value is not blank and is not more than 100 characters
	// long
	if strings.TrimSpace(form.Title) == "" {
		form.FieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(form.Title) > 100 {
		form.FieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	// Verifying if content isn't blank
	if strings.TrimSpace(form.Content) == "" {
		form.FieldErrors["content"] = "This field cannot be blank"
	}

	// Verifying if expires value matches one of the permitted values
	if form.Expires != 1 && form.Expires != 7 && form.Expires != 365 {
		form.FieldErrors["expires"] = "This field must be equal to 1, 7 or 365"
	}

	// If there is any errors, re-display the "create.tmpl.html" passing in the
	// snippetCreateForm instance to dynamic data in the Form field
	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData(r)
		data.Form = form

		app.render(w, 422, "create.tmpl.html", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), 303)
}

// Handler to see the form that will be use to create a snippet
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	// Initializing a new snippetCreateForm instance and passing it to the template
	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, 200, "create.tmpl.html", data)
}
