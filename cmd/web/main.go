package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"text/template"

	_ "github.com/go-sql-driver/mysql"

	// Database Snippet Model
	"lucasvinibox.isaacszf.net/internal/models"
)

// Application struct to hold the application-wide dependencies for the application.
// (it will be used for dependency injection)
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

// Function to open the database and verifying it
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Verifying if connection succeeded
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	// Adding a new command-line flag to set the port (default is "4000")
	// EX: go ./cmd/web -addr=":9999"
	addr := flag.String("addr", ":4000", "HTTP Network Address")
	dsn := flag.String("dsn",
		"web:220406@/lucasvinibox?parseTime=true", "MySQL data source name")

	// This is needed because this reads the command file passed by the user and parse it
	// to the variable "addr". Otherwise, "addr" will always be ":4000"
	flag.Parse()

	// Creating a new Logger for error and info messages.
	// Parameters -> Destination (os.Stdout = standard input), String Prefix for the msg,
	// some additional information to include
	// (flags are joined using the "|" operator)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Creating a connection pool between the App and Database
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Initializing a new template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	// Initializing App Struct and ServeMux
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &models.SnippetModel{Database: db},
		templateCache: templateCache,
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
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
