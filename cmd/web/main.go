package main

import (
	"log"
	"net/http"
)

func main() {
	// Initializing a new servemux, then register the "home" function
	// as the handler for the "/" URL pattern
	mux := http.NewServeMux()

	/* I still not sure why a public file server is useful in this project but OK */

	// Creating file server which get files from "ui/static" to the user
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// Using the "mux.Handle()" to register the file server as the handler for all
	// URL paths that start with "/static/". We need to use "StripPrefix" to match
	// with the dir path
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Println("🔥 Starting server on :4000")

	// Starting web server with "listenAndAdvice"
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
