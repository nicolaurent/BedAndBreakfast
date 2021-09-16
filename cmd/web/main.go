package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/nicolaurent/bedandbreakfast/internal/config"
	"github.com/nicolaurent/bedandbreakfast/internal/driver"
	"github.com/nicolaurent/bedandbreakfast/internal/handlers"
	"github.com/nicolaurent/bedandbreakfast/internal/helpers"
	"github.com/nicolaurent/bedandbreakfast/internal/models"
	"github.com/nicolaurent/bedandbreakfast/internal/renders"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main applicatio function
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	// what am I going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	// Change this to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to database
	log.Println("Connection to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=booking user=postgres password=password")

	if err != nil {
		log.Fatal("cannot connect to database! Dying...")
	}

	tc, err := renders.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	renders.NewRenderer(&app)
	helpers.NewHelpers(&app)

	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)

	fmt.Println("Starting application on port", portNumber)

	// _ = http.ListenAndServe(portNumber, nil)

	return db, nil
}
