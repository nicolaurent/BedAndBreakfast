package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nicolaurent/bedandbreakfast/pkg/config"
	"github.com/nicolaurent/bedandbreakfast/pkg/handlers"
	"github.com/nicolaurent/bedandbreakfast/pkg/renders"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

// main is the main applicatio function
func main() {
	/*
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			n, err := fmt.Fprintf(w, "Hello, world!")
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Bytes written:", n)
		})
		// */

	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := renders.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	renders.NewTemplates(&app)

	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)

	fmt.Println("Starting application on port", portNumber)

	// _ = http.ListenAndServe(portNumber, nil)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}