package main

import (
	"net/http"

	"github.com/justinas/alice"
)

// Creates a ServeMux and set all routes
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	/* I still not sure why a public file server is useful in this project but OK */

	// Creating file server which get files from "ui/static" to the user
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// Using the "mux.Handle()" to register the file server as the handler for all
	// URL paths that start with "/static/". We need to use "StripPrefix" to match
	// with the dir path
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	// Passing all the middlewares before mux
	// This is needed so that the middleware can acts in all requests
	// BEFORE: "app.recoverPanic(app.logRequest(secureHeaders(mux)))"

	// Creating a middleware chain containing our 'standard' middleware
	// which will be used for every request our app receives
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(mux)
}
