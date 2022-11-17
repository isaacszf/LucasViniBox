package main

import "net/http"

// Creates a ServeMux and set all routes
func (app *application) routes() *http.ServeMux {
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

	return mux
}
