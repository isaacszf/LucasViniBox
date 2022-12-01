package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// Creates a ServeMux and set all routes
func (app *application) routes() http.Handler {
	// Using custom router (httprouter)
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	/* I still not sure why a public file server is useful in this project but OK */

	// Creating file server which get files from "ui/static" to the user
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// Using the "router.Handler()" to register the file server as the handler for all
	// URL paths that start with "/static/". We need to use "StripPrefix" to match
	// with the dir path
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreate)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)

	// Passing all the middlewares before router
	// This is needed so that the middleware can acts in all requests
	// BEFORE: "app.recoverPanic(app.logRequest(secureHeaders(router)))"

	// Creating a middleware chain containing our 'standard' middleware
	// which will be used for every request our app receives
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}
