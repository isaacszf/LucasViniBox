package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// Application struct to hold the application-wide dependencies for the application.
// (it will be used for dependency injection)
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	// Adding a new command-line flag to set the port (default is "4000")
	// EX: go ./cmd/web -addr=":9999"
	addr := flag.String("addr", ":4000", "HTTP Network Address")

	// This is needed because this reads the command file passed by the user and parse it
	// to the variable "addr". Otherwise, "addr" will always be ":4000"
	flag.Parse()

	// Creating a new Logger for error and info messages.
	// Parameters -> Destination (os.Stdout = standard input), String Prefix for the msg,
	// some additional information to include
	// (flags are joined using the "|" operator)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Initializing App Struct and ServeMux
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	muxRoutes := app.routes()

	// Initializing a new http.Server struct. We set the Addr and Handler fields so
	// that the server uses the same network address and routes as before, and set
	// the ErrorLog field so that the server now uses the custom errorLog logger in
	// the event of any problems.
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  muxRoutes,
	}

	infoLog.Printf("🔥 Starting server on http://localhost%s", *addr)

	// Starting web server with "listenAndAdvice" using the custom struct
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
